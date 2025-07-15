package zoneserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type zoneServerSession struct {
	server   *Server
	conn     net.Conn
	id       uint32
	sendChan chan []byte
	done     chan struct{}
	agentId  byte
	wg       sync.WaitGroup
}

func newZoneServerSession(id uint32, conn net.Conn) network.TCPServerSession {
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		_ = tcpConn.SetNoDelay(true)
	}

	session := &zoneServerSession{
		id:       id,
		conn:     conn,
		sendChan: make(chan []byte, 100),
		done:     make(chan struct{}),
	}

	session.wg.Add(1)
	go session.sender()

	return session
}

func (s *zoneServerSession) ID() uint32 {
	return s.id
}

func (s *zoneServerSession) Handle() {
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

func (s *zoneServerSession) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (s *zoneServerSession) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *zoneServerSession) processPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	ctrl := packet[8]
	cmd := packet[9]
	var proto uint16
	if len(packet) > 11 {
		proto = binary.LittleEndian.Uint16(packet[10:])
	}

	s.server.Logger.Error(
		"Unhandled packet",
		shared.Field{Key: "ctrl", Value: ctrl},
		shared.Field{Key: "cmd", Value: cmd},
		shared.Field{Key: "protocolo", Value: proto},
	)
}

func (s *zoneServerSession) sender() {
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
