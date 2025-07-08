package mainserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type mainServerSession struct {
	server   *Server
	conn     net.Conn
	id       uint32
	sendChan chan []byte
	done     chan struct{}
	serverId byte
	wg       sync.WaitGroup
}

func newMainServerSession(id uint32, conn net.Conn) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &mainServerSession{
		id:       id,
		conn:     conn,
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *mainServerSession) ID() uint32 {
	return s.id
}

func (s *mainServerSession) Handle() {
	defer func() {
		s.server.Logger.Info(fmt.Sprintf("Server %d disconnected", s.serverId))
		s.server.RemoveSession(s.id)
		close(s.done)
		s.wg.Wait()
	}()
	buffer := make([]byte, 1024*16)
	dynamicBuffer := bytes.NewBuffer(nil)
	for {
		n, err := s.conn.Read(buffer)
		if err != nil {
			break
		}

		if n == 5 {
			packet := buffer[:n]
			s.serverId = packet[4]
			s.server.Logger.Info(fmt.Sprintf("Server %d connected", s.serverId),
				shared.Field{Key: "serverId", Value: s.serverId},
			)
			continue
		}

		dynamicBuffer.Write(buffer[:n])
		for dynamicBuffer.Len() >= 4 {
			dataLength := int(binary.LittleEndian.Uint16(dynamicBuffer.Bytes()[2:]))
			if dataLength > dynamicBuffer.Len() || dataLength == 0 {
				break
			}

			currentPacket := dynamicBuffer.Next(dataLength)
			go s.processPacket(currentPacket)
		}
	}
}

func (s *mainServerSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *mainServerSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *mainServerSession) processPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	proto := binary.LittleEndian.Uint16(packet)
	switch proto {
	case protocol.S2MCharacterLogin:
		msg, err := messages.ReadMsgS2MCharacterLogin(packet)
		if err != nil {
			return
		}

		account := utils.ReadStringFromBytes(msg.Account[:])
		characterName := utils.ReadStringFromBytes(msg.CharacterName[:])
		clientIp := utils.ReadStringFromBytes(msg.ClientIp[:])
		if s.server.players.HasPlayer(msg.PcId) {
			errMsg := messages.NewMsgM2SError(msg.PcId, constants.ErrorCodeCharacterLoginFailed, "Player already logged in", msg.GateServerId)
			_ = s.Send(errMsg.GetBytes())
			s.server.Logger.Info("Player already logged in",
				shared.Field{Key: "pcId", Value: msg.PcId},
				shared.Field{Key: "serverId", Value: s.serverId},
				shared.Field{Key: "gateServerId", Value: msg.GateServerId},
				shared.Field{Key: "characterName", Value: characterName},
				shared.Field{Key: "account", Value: account},
			)
			return
		}

		mapId, err := s.server.dbService.GetCharacterMapInfo(msg.PcId, characterName)
		if err != nil {
			errMsg := messages.NewMsgM2SError(msg.PcId, constants.ErrorCodeCharacterLoginFailed, "Character not found", msg.GateServerId)
			_ = s.Send(errMsg.GetBytes())
			s.server.Logger.Error("Failed to get character map info",
				shared.Field{Key: "error", Value: err},
				shared.Field{Key: "serverId", Value: s.serverId},
				shared.Field{Key: "gateServerId", Value: msg.GateServerId},
				shared.Field{Key: "characterName", Value: characterName},
				shared.Field{Key: "account", Value: account},
			)
			return
		}

		zone, exists := s.server.mapZones.Get(mapId)
		if !exists {
			errMsg := messages.NewMsgM2SError(msg.PcId, constants.ErrorCodeCharacterLoginFailed, "Character zone not found", msg.GateServerId)
			_ = s.Send(errMsg.GetBytes())
			s.server.Logger.Error("Character zone not found",
				shared.Field{Key: "mapId", Value: mapId},
				shared.Field{Key: "serverId", Value: s.serverId},
				shared.Field{Key: "gateServerId", Value: msg.GateServerId},
				shared.Field{Key: "characterName", Value: characterName},
				shared.Field{Key: "account", Value: account},
			)
			return
		}

		player := NewPlayer(
			msg.PcId,
			account,
			characterName,
			clientIp,
			mapId,
			zone.serverId,
			msg.GateServerId,
			zone,
		)
		s.server.players.Add(player)
		loginMsg := messages.NewMsgM2SAnsCharacterLogin(msg.PcId, msg.GateServerId, mapId, 0)
		_ = s.Send(loginMsg.GetBytes())
	case protocol.S2MMapList:
		mapCount := binary.LittleEndian.Uint16(packet[10:])
		mapIds := make([]uint16, mapCount)
		for i := 0; i < int(mapCount); i++ {
			mapIds[i] = binary.LittleEndian.Uint16(packet[12+i*2:])
		}

		if s.server.IsZoneRegistered(s.serverId) {
			return
		}

		zone := NewZone(s.serverId, s)
		zone.SetMaps(mapIds)
		s.server.zoneSessions.Set(s.serverId, zone)
		for _, mapId := range mapIds {
			s.server.mapZones.Set(mapId, zone)
		}

	default:
		s.server.Logger.Info("Unhandled packet",
			shared.Field{Key: "packet", Value: packet},
			shared.Field{Key: "serverId", Value: s.serverId})
	}
}

func (s *mainServerSession) sender() {
	defer s.wg.Done()
	for {
		select {
		case data := <-s.sendChan:
			if _, err := s.conn.Write(data); err != nil {
				s.server.Logger.Error(
					fmt.Sprintf("Failed to send packet to server %d", s.serverId),
					shared.Field{Key: "error", Value: err},
					shared.Field{Key: "sessionId", Value: s.id},
					shared.Field{Key: "serverId", Value: s.serverId})
				return
			}

		case <-s.done:
			return
		}
	}
}
