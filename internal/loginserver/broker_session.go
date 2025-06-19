package loginserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
)

type brokerSession struct {
	server        *Broker
	conn          net.Conn
	id            uint32
	serverId      byte
	serverName    string
	port          uint32
	ipAddress     string
	isInitialized bool
	sendChan      chan []byte
	done          chan struct{}
	wg            sync.WaitGroup
}

func newBrokerSession(id uint32, conn net.Conn, server interface{}) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &brokerSession{
		id:       id,
		conn:     conn,
		server:   server.(*Broker),
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *brokerSession) ID() uint32 {
	return s.id
}

func (s *brokerSession) Handle() {
	defer func() {
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

		s.processPacket(packet)
	}
}

func (s *brokerSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *brokerSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *brokerSession) processPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	ctrl := packet[8]
	cmd := packet[9]
	switch ctrl {
	case 0x02:
		switch cmd {
		case 0xE0:
			s.handleGateConnect(packet)
		case 0xE2:
			s.handleAccountLogout(packet)
		case 0xE3:
			s.handleAccountLogin(packet)
		}
	}
}

func (s *brokerSession) handleGateConnect(packet []byte) {
	if s.isInitialized {
		return
	}

	msg, err := messages.ReadMsgGate2LsConnect(packet)
	if err != nil {
		s.server.Logger.Error("Failed to read gate connect message", shared.Field{Key: "error", Value: err})
		return
	}

	s.serverId = msg.ServerId
	s.ipAddress = utils.ReadStringFromBytes(msg.IpAddress[:])
	s.port = msg.Port
	s.serverName = utils.ReadStringFromBytes(msg.Name[:])
	s.isInitialized = true
	s.server.Logger.Info(
		fmt.Sprintf("Gate Server %d connected", s.serverId),
		shared.Field{Key: "ipAddress", Value: s.ipAddress},
		shared.Field{Key: "port", Value: s.port},
		shared.Field{Key: "serverName", Value: s.serverName},
		shared.Field{Key: "id", Value: s.id},
	)
}

func (s *brokerSession) handleAccountLogout(packet []byte) {
	msg, err := messages.ReadMsgGate2LsAccLogout(packet)
	if err != nil {
		s.server.Logger.Error("Failed to read gate account logout message", shared.Field{Key: "error", Value: err})
		return
	}

	account := utils.ReadStringFromBytes(msg.Account[:])
	s.server.Logger.Info(fmt.Sprintf("%s logged out", account))
	s.server.RemoveLoggedInUser(account)
}

func (s *brokerSession) handleAccountLogin(packet []byte) {
	msg, err := messages.ReadMsgGate2LsPreparedAccLogin(packet)
	if err != nil {
		s.server.Logger.Error("Failed to read gate account login message", shared.Field{Key: "error", Value: err})
		return
	}

	account := utils.ReadStringFromBytes(msg.Account[:])
	s.server.Logger.Info(fmt.Sprintf("%s logged in", account))
	s.server.AddLoggedInUser(account, msg.PcId)
}

func (s *brokerSession) sender() {
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
