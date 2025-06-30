package accountserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type accountServerSession struct {
	server   *Server
	conn     net.Conn
	id       uint32
	sendChan chan []byte
	done     chan struct{}
	agentId  byte
	wg       sync.WaitGroup
	players  *shared.SafeMap[uint32, *Player]
}

func newAccountServerSession(id uint32, conn net.Conn) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &accountServerSession{
		id:       id,
		conn:     conn,
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
		players:  shared.NewSafeMap[uint32, *Player](),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *accountServerSession) ID() uint32 {
	return s.id
}

func (s *accountServerSession) Handle() {
	defer func() {
		s.server.Logger.Info(fmt.Sprintf("Gate server %d disconnected", s.agentId))
		s.server.RemoveSession(s.id)
		close(s.done)
		s.wg.Wait()
	}()
	for {
		var buf bytes.Buffer
		if _, err := io.CopyN(&buf, s.conn, 4); err != nil {
			break
		}

		reader := io.MultiReader(&buf, s.conn)
		dataLength := binary.LittleEndian.Uint32(buf.Bytes())
		if dataLength == 0 {
			continue
		}

		if dataLength > 16*1024*1024 {
			break
		}

		packet := make([]byte, dataLength)
		if _, err := io.ReadFull(reader, packet); err != nil {
			break
		}

		go s.processPacket(packet)
	}
}

func (s *accountServerSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *accountServerSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *accountServerSession) processPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	ctrl := packet[8]
	cmd := packet[9]
	switch ctrl {
	case 0x01:
		switch cmd {
		case 0xE0:
			s.handleGateConnect(packet)
		case 0xE1:
			s.handleCharacterListing(packet)
		case 0xE2:
			s.handleClientDisconnect(packet)
		default:
			s.server.Logger.Error("Unhandled packet", shared.Field{Key: "ctrl", Value: ctrl}, shared.Field{Key: "cmd", Value: cmd})
		}

	case 0x03:
		switch cmd {
		case 0xFF:
			s.handleProtocolPacket(packet)
		default:
			s.server.Logger.Error("Unhandled packet", shared.Field{Key: "ctrl", Value: ctrl}, shared.Field{Key: "cmd", Value: cmd})
		}

	default:
		s.server.Logger.Error("Unhandled packet", shared.Field{Key: "ctrl", Value: ctrl}, shared.Field{Key: "cmd", Value: cmd})
	}
}

func (s *accountServerSession) handleGateConnect(packet []byte) {
	s.server.Logger.Info(fmt.Sprintf("Gate server %d connected", packet[10]))
	s.agentId = packet[10]
}

func (s *accountServerSession) handleCharacterListing(packet []byte) {
	pcId := binary.LittleEndian.Uint32(packet[4:])
	_, exists := s.players.Get(pcId)
	if exists || pcId == 0 {
		_ = s.sendErrorMsg(pcId, constants.ErrorCodeLoginFailed, constants.AccountAlreadyLoggedInMsg)
		return
	}

	msg, err := messages.ReadMsgGate2AsNewClient(packet)
	if err != nil {
		_ = s.sendErrorMsg(pcId, constants.ErrorCodeLoginFailed, constants.LoginFailedMsg)
		return
	}

	player := NewPlayer(pcId, utils.ReadStringFromBytes(msg.Account[:]), utils.ReadStringFromBytes(msg.ClientIP[:]))
	s.players.Set(pcId, player)
	characters, err := s.server.dbService.GetCharactersForListing(pcId)
	if err != nil {
		_ = s.sendErrorMsg(pcId, constants.ErrorCodeLoginFailed, constants.LoginFailedMsg)
		return
	}

	if len(characters) == 0 {
		msg := messages.NewMsgS2CCharacterListEmpty(pcId)
		data := msg.GetBytes()
		_ = s.Send(data)
		return
	}

	characterList := make([]messages.CharacterInfo, len(characters))
	for i, character := range characters {
		lastUsed := byte(0)
		if i == 0 {
			lastUsed = 1
		}

		characterList[i] = messages.CharacterInfo{
			LastUsed: lastUsed,
			Class:    character.Class,
			Level:    character.Level,
			Nation:   character.Data.SocialInfo.Nation,
		}
		copy(characterList[i].Name[:], utils.MakeFixedLengthStringBytes(character.Name, 0x15))
		for j := 0; j < len(character.Data.Wear); j++ {
			if j > 9 {
				break
			}

			item, exists := s.server.GetItem(character.Data.Wear[j].ItemCode)
			if !exists {
				continue
			}

			characterList[i].Wear[j] = messages.CharacterWear{
				ItemPtr:    0,
				ItemCode:   character.Data.Wear[j].ItemCode,
				ItemOption: character.Data.Wear[j].ItemOption,
				WearIndex:  uint32(item.SlotIndex),
			}
		}
	}

	listMsg := messages.NewMsgS2CCharacterList(pcId, characterList)
	data := listMsg.GetBytes()
	_ = s.Send(data)
}

func (s *accountServerSession) handleClientDisconnect(packet []byte) {
	if len(packet) < 8 {
		return
	}

	pcId := binary.LittleEndian.Uint32(packet[4:])
	player, exists := s.players.Get(pcId)
	if !exists {
		return
	}

	s.server.Logger.Info(fmt.Sprintf("Account %s disconnected", player.account))
	s.players.Delete(pcId)
}

func (s *accountServerSession) sender() {
	defer s.wg.Done()
	for {
		select {
		case data := <-s.sendChan:
			if _, err := s.conn.Write(data); err != nil {
				s.server.Logger.Error("Failed to send packet to gate server",
					shared.Field{Key: "error", Value: err},
					shared.Field{Key: "sessionId", Value: s.id})
				return
			}

		case <-s.done:
			return
		}
	}
}

func (s *accountServerSession) sendErrorMsg(pcId uint32, errorCode uint16, errorMsg string) error {
	msg := messages.NewMsgS2CError(pcId, errorCode, errorMsg)
	data := msg.GetBytes()
	return s.Send(data)
}

func (s *accountServerSession) handleProtocolPacket(packet []byte) {
	if len(packet) < 12 {
		return
	}

	proto := binary.LittleEndian.Uint16(packet[10:])
	switch proto {
	case protocol.C2SCharLogout:
		s.handleClientDisconnect(packet)
	default:
		s.server.Logger.Error("Unhandled packet", shared.Field{Key: "protocol", Value: proto})
	}
}
