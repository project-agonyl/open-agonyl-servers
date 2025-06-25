package gateserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	sharedconstants "github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type serverSession struct {
	server   *Server
	conn     net.Conn
	id       uint32
	player   *Player
	sendChan chan []byte
	done     chan struct{}
	wg       sync.WaitGroup
}

func newServerSession(id uint32, conn net.Conn) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &serverSession{
		id:       id,
		conn:     conn,
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *serverSession) ID() uint32 {
	return s.id
}

func (s *serverSession) Handle() {
	defer func() {
		close(s.done)
		s.wg.Wait()
		s.sendServerLogoutMsg()
		if s.player != nil {
			s.server.players.Remove(s.player.Id)
		}

		s.server.RemoveSession(s.id)
		_ = s.server.db.SetAccountOffline(s.player.Id)
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

		s.processPacket(packet)
	}
}

func (s *serverSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *serverSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *serverSession) processPacket(packet []byte) {
	ctrl := packet[8]
	cmd := packet[9]
	var protocol uint16
	if len(packet) > 11 {
		protocol = binary.LittleEndian.Uint16(packet[10:])
	}

	if s.player != nil {
		binary.LittleEndian.PutUint32(packet[4:], s.player.Id)
	}

	switch ctrl {
	case 0x01:
		switch cmd {
		case 0xE2:
			s.handleLogin(packet)
		case 0xF0:
			// TODO: Implement handling of ping packet
		}

	case 0x03:
		if s.player == nil {
			return
		}

		s.server.crypto.Decrypt(packet)
		switch protocol {
		case 0x1106: // Character login
			fallthrough
		case 0x2322: // Transfer Clan Mark
			fallthrough
		case 0x2323: // Clan...
			fallthrough
		case 0xA001: // Create Character
			fallthrough
		case 0xA002: // Delete Character
			_ = s.server.zoneServerClients.Send(constants.AccountServerServerId, packet)
		default:
			_ = s.server.zoneServerClients.Send(s.player.GetCurrentZone(), packet)
		}

	default:
		if s.player == nil {
			return
		}

		s.server.crypto.Decrypt(packet)
		switch protocol {
		case 0x2322: // Transfer Clan Mark
			fallthrough
		case 0x2323: // Clan...
			_ = s.server.zoneServerClients.Send(constants.AccountServerServerId, packet)
		}
	}
}

func (s *serverSession) sender() {
	defer s.wg.Done()
	for {
		select {
		case data := <-s.sendChan:
			if _, err := s.conn.Write(data); err != nil {
				s.server.Logger.Error("Failed to send packet to client",
					shared.Field{Key: "error", Value: err},
					shared.Field{Key: "sessionId", Value: s.id})
				return
			}

		case <-s.done:
			return
		}
	}
}

func (s *serverSession) handleLogin(packet []byte) {
	msg, err := messages.ReadMsgC2SGateLogin(packet)
	if err != nil {
		s.server.Logger.Error("Failed to read login message",
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "sessionId", Value: s.id})
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.LoginFailedMsg)
		return
	}

	id := msg.PcId
	username := utils.ReadStringFromBytes(msg.Account[:])
	password := utils.ReadStringFromBytes(msg.Password[:])
	if !s.server.loginServerClient.IsLoggedIn(id) {
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.LoginFailedMsg)
		return
	}

	if s.server.players.HasPlayer(id) {
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.AccountAlreadyLoggedInMsg)
		return
	}

	account, err := s.server.db.GetAccount(id)
	if err != nil {
		s.server.Logger.Error("Failed to get account",
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "sessionId", Value: s.id})
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.LoginFailedMsg)
		return
	}

	if account.Username != username {
		s.server.Logger.Error("Account username mismatch",
			shared.Field{Key: "expected", Value: account.Username},
			shared.Field{Key: "actual", Value: username},
			shared.Field{Key: "sessionId", Value: s.id})
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.LoginFailedMsg)
		return
	}

	if account.IsOnline {
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.AccountAlreadyLoggedInMsg)
		return
	}

	if err := s.server.db.SetAccountOnline(id, account); err != nil {
		s.server.Logger.Error("Failed to set account online",
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "sessionId", Value: s.id})
		_ = s.sendErrorMsg(sharedconstants.ErrorCodeLoginFailed, sharedconstants.LoginFailedMsg)
		return
	}

	s.player = NewPlayer(id, username, s.conn, s.server.Logger)
	go func(session *serverSession) {
		_ = session.server.loginServerClient.Send(messages.NewMsgGate2LsPreparedAccLogin(username).GetBytes())
	}(s)
	s.server.Logger.Info(
		fmt.Sprintf("Account %s session started", username),
		shared.Field{Key: "id", Value: id},
		shared.Field{Key: "username", Value: username},
	)
	_ = s.server.zoneServerClients.Send(
		constants.AccountServerServerId,
		messages.NewMsgGate2AsNewClient(
			username, password, strings.Split(s.conn.RemoteAddr().String(), ":")[0], id,
		).GetBytes(),
	)
}

func (s *serverSession) sendErrorMsg(errorCode uint16, errorMsg string) error {
	msg := messages.NewMsgS2CError(0, errorCode, errorMsg)
	data := msg.GetBytes()
	s.server.crypto.Encrypt(data)
	_, err := s.conn.Write(data)
	return err
}

func (s *serverSession) sendServerLogoutMsg() {
	if s.player == nil {
		return
	}

	s.server.Logger.Info(
		fmt.Sprintf("Account %s session ended", s.player.Username),
		shared.Field{Key: "id", Value: s.player.Id},
		shared.Field{Key: "username", Value: s.player.Username},
	)
	_ = s.server.zoneServerClients.Send(s.player.GetCurrentZone(), messages.NewMsgZa2ZsAccLogout(s.player.Id, 0x00).GetBytes())
	_ = s.server.loginServerClient.Send(messages.NewMsgGate2LsAccLogout(0x00, s.player.Username).GetBytes())
}
