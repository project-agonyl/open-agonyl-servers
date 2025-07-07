package mainserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
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
	if len(packet) < 5 {
		return
	}

	s.server.Logger.Info("Unhandled packet",
		shared.Field{Key: "packet", Value: packet},
		shared.Field{Key: "sessionId", Value: s.id},
		shared.Field{Key: "serverId", Value: s.serverId})
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
