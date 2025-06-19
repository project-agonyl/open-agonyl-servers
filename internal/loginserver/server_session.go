package loginserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
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
		tcpConn.SetNoDelay(true)
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
	msg, err := messages.ReadMsgC2SLogin(packet)
	if err != nil {
		s.sendClientMsg(constants.InvalidCredentialsMsg)
		s.server.Logger.Error("Could not read login message",
			shared.Field{Key: "error", Value: err})
		return
	}

	_ = utils.ReadStringFromBytes(msg.Username[:])
	_ = utils.ReadStringFromBytes(msg.Password[:])
}

func (s *loginServerSession) handleServerSelect(serverId byte) {}
