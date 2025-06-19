package loginserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type loginServerSession struct {
	server   *Server
	conn     net.Conn
	id       uint32
	account  *db.Account
	sendChan chan []byte
	done     chan struct{}
	wg       sync.WaitGroup
}

func newLoginServerSession(id uint32, conn net.Conn, server interface{}) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &loginServerSession{
		id:       id,
		conn:     conn,
		server:   server.(*Server),
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *loginServerSession) ID() uint32 {
	return s.id
}

func (s *loginServerSession) Handle() {
	defer func() {
		s.server.RemoveSession(s.id)
		if s.account != nil {
			s.server.RemoveLoggedInUser(s.account.Username)
		}

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

		s.processPacket(packet)
	}
}

func (s *loginServerSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *loginServerSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *loginServerSession) processPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	if s.server.broker.GetGateServerCount() == 0 {
		_ = s.sendClientMsg(constants.ServerUnderMaintenanceMsg)
		return
	}

	ctrl := packet[8]
	cmd := packet[9]
	switch ctrl {
	case 0x01:
		switch cmd {
		case 0xE0:
			s.handleLogin(packet)
		case 0xE1:
			s.handleServerSelect(packet[10])
		}
	}
}

func (s *loginServerSession) sender() {
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

func (s *loginServerSession) sendClientMsg(msg string) error {
	msgPacket := messages.NewMsgLs2ClSay(msg)
	return binary.Write(s.conn, binary.LittleEndian, msgPacket)
}

func (s *loginServerSession) handleLogin(packet []byte) {
	if s.account != nil {
		return
	}

	msg, err := messages.ReadMsgC2SLogin(packet)
	if err != nil {
		_ = s.sendClientMsg(constants.InvalidCredentialsMsg)
		s.server.Logger.Error("Could not read login message",
			shared.Field{Key: "error", Value: err})
		return
	}

	username := strings.TrimSpace(utils.ReadStringFromBytes(msg.Username[:]))
	password := strings.TrimSpace(utils.ReadStringFromBytes(msg.Password[:]))
	if username == "" || password == "" {
		_ = s.sendClientMsg(constants.InvalidCredentialsMsg)
		return
	}

	account, err := s.server.dbService.GetAccountByUsername(username)
	if err != nil {
		_ = s.sendClientMsg(constants.InvalidCredentialsMsg)
		s.server.Logger.Error("Could not get account by username",
			shared.Field{Key: "error", Value: err})
		return
	}

	if !s.server.cfg.IsTestMode {
		err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
		if err != nil {
			_ = s.sendClientMsg(constants.InvalidCredentialsMsg)
			return
		}

		if strings.EqualFold(account.AccountStatus, constants.AccountStatusBanned) {
			_ = s.sendClientMsg(constants.AccountBannedMsg)
			return
		}

		if !strings.EqualFold(account.AccountStatus, constants.AccountStatusActive) {
			_ = s.sendClientMsg(constants.AccountNotActiveMsg)
			return
		}
	}

	if account.IsOnline || s.server.IsLoggedIn(username) {
		_ = s.sendClientMsg(constants.AccountAlreadyLoggedInMsg)
		return
	}

	s.account = account
	header := messages.MsgHeadNoProtocol{
		Ctrl: 0x01,
		Cmd:  0xE1,
		PcId: s.account.ID,
	}
	var serverInfoBuffer bytes.Buffer
	_ = binary.Write(&serverInfoBuffer, binary.LittleEndian, &header)
	_ = binary.Write(&serverInfoBuffer, binary.LittleEndian, uint16(s.server.broker.GetGateServerCount()))
	for serverId, serverName := range s.server.broker.GetGateServerList() {
		serverInfo := messages.GateServerInfo{
			ServerID: serverId,
		}
		copy(serverInfo.ServerName[:], utils.MakeFixedLengthStringBytes(serverName, 0x11))
		copy(serverInfo.ServerStatus[:], utils.MakeFixedLengthStringBytes("ONLINE", 0x51))
		_ = binary.Write(&serverInfoBuffer, binary.LittleEndian, &serverInfo)
	}

	serverInfoBufferBytes := serverInfoBuffer.Bytes()
	binary.LittleEndian.PutUint32(serverInfoBufferBytes, uint32(len(serverInfoBufferBytes)))
	s.server.AddLoggedInUser(username, s.account.ID)
	_ = s.Send(serverInfoBufferBytes)
	s.server.Logger.Info(
		fmt.Sprintf("Account %s logged in", username),
		shared.Field{Key: "username", Value: username},
		shared.Field{Key: "id", Value: s.account.ID},
	)
}

func (s *loginServerSession) handleServerSelect(serverId byte) {
	if s.account == nil {
		return
	}

	serverInfo, err := s.server.broker.GetGateServerInfoByServerId(serverId)
	if err != nil {
		_ = s.sendClientMsg(constants.ServerUnderMaintenanceMsg)
		s.server.Logger.Error(
			fmt.Sprintf("Failed to get gate server info by server id %d", serverId),
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "serverId", Value: serverId},
			shared.Field{Key: "accountId", Value: s.account.ID},
			shared.Field{Key: "username", Value: s.account.Username},
		)
		return
	}

	gateServerMsg := messages.NewMsgLs2GateLogin(s.account.Username, s.account.ID)
	if err := s.server.broker.SendMsgToGateServer(serverInfo.Id, gateServerMsg.GetBytes()); err != nil {
		_ = s.sendClientMsg(constants.ServerUnderMaintenanceMsg)
		s.server.Logger.Error(
			fmt.Sprintf("Failed to send message to gate server %d", serverId),
			shared.Field{Key: "error", Value: err},
			shared.Field{Key: "serverId", Value: serverInfo.Id},
			shared.Field{Key: "accountId", Value: s.account.ID},
			shared.Field{Key: "username", Value: s.account.Username},
		)
		return
	}

	gateServerInfoMsg := messages.NewMsgS2CGateInfo(s.account.ID, serverInfo.IpAddress, serverInfo.Port)
	_ = binary.Write(s.conn, binary.LittleEndian, gateServerInfoMsg)
}
