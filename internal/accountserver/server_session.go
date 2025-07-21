package accountserver

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/db"
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
	if pcId == 0 {
		_ = s.sendErrorMsg(pcId, constants.ErrorCodeLoginFailed, constants.AccountAlreadyLoggedInMsg)
		return
	}

	msg, err := messages.ReadMsgGate2AsNewClient(packet)
	if err != nil {
		_ = s.sendErrorMsg(pcId, constants.ErrorCodeLoginFailed, constants.LoginFailedMsg)
		return
	}

	player := NewPlayer(pcId, utils.ReadStringFromBytes(msg.Account[:]), utils.ReadStringFromBytes(msg.ClientIP[:]), s)
	s.server.players.Add(player)
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
		characterList[i] = messages.CharacterInfo{
			SlotUsed: 1,
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

			characterList[i].Wear[j] = messages.AclCharacterWear{
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
	pcId := binary.LittleEndian.Uint32(packet[4:])
	player, exists := s.server.players.Get(pcId)
	if !exists {
		return
	}

	if player.selectedCharacterName != "" {
		msg := messages.NewMsgS2MCharacterLogout(pcId, player.selectedCharacterName)
		_ = s.server.mainServerClient.Send(msg.GetBytes())
	}

	s.server.Logger.Info(fmt.Sprintf("Account %s disconnected", player.account))
	s.server.players.Remove(pcId)
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
	case protocol.C2SCharacterLogout:
		s.handleClientDisconnect(packet)
	case protocol.C2SAskCreatePlayer:
		s.handleCharacterCreate(packet)
	case protocol.C2SAskDeletePlayer:
		s.handleCharacterDelete(packet)
	case protocol.C2SCharacterLogin:
		s.handleCharacterLogin(packet)
	default:
		s.server.Logger.Error("Unhandled packet from gate server", shared.Field{Key: "protocol", Value: proto})
	}
}

func (s *accountServerSession) handleCharacterCreate(packet []byte) {
	msg, err := messages.ReadMsgC2SAskCreatePlayer(packet)
	if err != nil {
		return
	}

	name := utils.ReadStringFromBytes(msg.Name[:])
	exists, err := s.server.dbService.DoesCharacterExist(name)
	if exists || err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeDuplicateCharacter, constants.DuplicateCharacterMsg)
		if err != nil {
			s.server.Logger.Error(
				"Failed to check if character exists",
				shared.Field{Key: "error", Value: err},
				shared.Field{Key: "pcId", Value: msg.PcId},
			)
		}

		return
	}

	count, err := s.server.dbService.GetCharacterCount(msg.PcId)
	if count >= constants.MaxCharactersPerAccount || err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.MaxCharactersPerAccountExceededMsg)
		s.server.Logger.Error(
			"Failed to get character count",
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "pcId", Value: msg.PcId},
		)
		return
	}

	serials, err := s.server.serialNumberGenerator.GetNextSerials(context.Background(), 7)
	if err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.LoginFailedMsg)
		s.server.Logger.Error(
			"Failed to get next serials",
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "pcId", Value: msg.PcId},
		)
		return
	}

	var data db.CharacterData
	switch msg.Town {
	case 0x01:
		data.SocialInfo.Nation = 1
		data.Location.MapCode = 7
		data.Location.Position.X = 110
		data.Location.Position.Y = 110
	default:
		data.SocialInfo.Nation = 0
		data.Location.MapCode = 1
		data.Location.Position.X = 110
		data.Location.Position.Y = 110
	}

	switch msg.Class {
	case 0x01:
		data.Stats.Strength = 30
		data.Stats.Dexterity = 20
		data.Stats.Vitality = 25
		data.Stats.Mana = 20
		data.Stats.HP = 50
		data.Stats.MP = 30
		data.Stats.HPCapacity = 110
		data.Stats.MPCapacity = 40
		data.Wear = []db.WearItem{
			{ItemCode: 1048, ItemOption: 0, ItemUniqueCode: serials[0]},
			{ItemCode: 3322, ItemOption: 0, ItemUniqueCode: serials[1]},
			{ItemCode: 3307, ItemOption: 0, ItemUniqueCode: serials[2]},
			{ItemCode: 3297, ItemOption: 0, ItemUniqueCode: serials[3]},
			{ItemCode: 3302, ItemOption: 0, ItemUniqueCode: serials[4]},
			{ItemCode: 3317, ItemOption: 0, ItemUniqueCode: serials[5]},
			{ItemCode: 3312, ItemOption: 0, ItemUniqueCode: serials[6]},
		}
	case 0x02:
		data.Stats.Strength = 20
		data.Stats.Intelligence = 26
		data.Stats.Dexterity = 12
		data.Stats.Vitality = 20
		data.Stats.Mana = 40
		data.Stats.HP = 30
		data.Stats.MP = 80
		data.Stats.HPCapacity = 30
		data.Stats.MPCapacity = 120
		data.Wear = []db.WearItem{
			{ItemCode: 2066, ItemOption: 0, ItemUniqueCode: serials[0]},
			{ItemCode: 3337, ItemOption: 0, ItemUniqueCode: serials[1]},
			{ItemCode: 3327, ItemOption: 0, ItemUniqueCode: serials[2]},
			{ItemCode: 3332, ItemOption: 0, ItemUniqueCode: serials[3]},
			{ItemCode: 3347, ItemOption: 0, ItemUniqueCode: serials[4]},
			{ItemCode: 3342, ItemOption: 0, ItemUniqueCode: serials[5]},
		}
	case 0x03:
		data.Stats.Strength = 30
		data.Stats.Dexterity = 16
		data.Stats.Vitality = 25
		data.Stats.Mana = 25
		data.Stats.HP = 75
		data.Stats.MP = 37
		data.Stats.HPCapacity = 100
		data.Stats.MPCapacity = 50
		data.Wear = []db.WearItem{
			{ItemCode: 1110, ItemOption: 0, ItemUniqueCode: serials[0]},
			{ItemCode: 3677, ItemOption: 0, ItemUniqueCode: serials[1]},
			{ItemCode: 3657, ItemOption: 0, ItemUniqueCode: serials[2]},
			{ItemCode: 3667, ItemOption: 0, ItemUniqueCode: serials[3]},
			{ItemCode: 3697, ItemOption: 0, ItemUniqueCode: serials[4]},
			{ItemCode: 3687, ItemOption: 0, ItemUniqueCode: serials[5]},
		}
	default:
		data.Stats.Strength = 30
		data.Stats.Dexterity = 16
		data.Stats.Vitality = 30
		data.Stats.Mana = 20
		data.Stats.HP = 75
		data.Stats.MP = 20
		data.Stats.HPCapacity = 120
		data.Stats.MPCapacity = 30
		data.Wear = []db.WearItem{
			{ItemCode: 1024, ItemOption: 0, ItemUniqueCode: serials[0]},
			{ItemCode: 3282, ItemOption: 0, ItemUniqueCode: serials[1]},
			{ItemCode: 3272, ItemOption: 0, ItemUniqueCode: serials[2]},
			{ItemCode: 3277, ItemOption: 0, ItemUniqueCode: serials[3]},
			{ItemCode: 3292, ItemOption: 0, ItemUniqueCode: serials[4]},
			{ItemCode: 3287, ItemOption: 0, ItemUniqueCode: serials[5]},
		}
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.LoginFailedMsg)
		return
	}

	_, err = s.server.dbService.CreateCharacter(msg.PcId, name, msg.Class, bytes)
	if err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.LoginFailedMsg)
		return
	}

	wear := [0xA]messages.AclCharacterWear{}
	for i := 0; i < len(data.Wear); i++ {
		if i > 9 {
			break
		}

		item, exists := s.server.GetItem(data.Wear[i].ItemCode)
		if !exists {
			continue
		}

		wear[i] = messages.AclCharacterWear{
			ItemPtr:    0,
			ItemCode:   data.Wear[i].ItemCode,
			ItemOption: data.Wear[i].ItemOption,
			WearIndex:  uint32(item.SlotIndex),
		}
	}

	replyMsg := messages.NewMsgS2CAnsCreatePlayer(msg.PcId, msg.Class, name, wear)
	_ = s.Send(replyMsg.GetBytes())
}

func (s *accountServerSession) handleCharacterDelete(packet []byte) {
	msg, err := messages.ReadMsgC2SAskDeletePlayer(packet)
	if err != nil {
		return
	}

	name := utils.ReadStringFromBytes(msg.Name[:])
	err = s.server.dbService.DeleteCharacter(msg.PcId, name)
	if err != nil {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.CharacterNotFoundMsg)
		return
	}

	replyMsg := messages.NewMsgS2CAnsDeletePlayer(msg.PcId, name)
	_ = s.Send(replyMsg.GetBytes())
}

func (s *accountServerSession) handleCharacterLogin(packet []byte) {
	msg, err := messages.ReadMsgC2SCharacterLogin(packet)
	if err != nil {
		return
	}

	characterName := utils.ReadStringFromBytes(msg.CharacterName[:])
	player, exists := s.server.players.Get(msg.PcId)
	if !exists {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.LoginFailedMsg)
		return
	}

	if player.GetSelectedCharacterName() != "" {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeCharacterInvalid, constants.InvalidCharacterMsg)
		return
	}

	if exists, err := s.server.dbService.DoesCharacterExist(characterName); err != nil || !exists {
		_ = s.sendErrorMsg(msg.PcId, constants.ErrorCodeChracterNotFound, constants.LoginFailedMsg)
		return
	}

	player.SetSelectedCharacterName(characterName)
	msMsg := messages.NewMsgS2MCharacterLogin(msg.PcId, player.account, "", characterName, player.clientIp, s.agentId)
	_ = s.server.mainServerClient.Send(msMsg.GetBytes())
}
