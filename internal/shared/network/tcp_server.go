package network

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type NewSessionFunc func(id uint32, conn net.Conn, server interface{}) TCPServerSession

type TCPServer struct {
	Logger       shared.Logger
	Name         string
	Addr         string
	Listener     net.Listener
	Sessions     *shared.SafeMap[uint32, TCPServerSession]
	Running      atomic.Bool
	NewSession   NewSessionFunc
	UidGenerator *shared.UidGenerator
}

func (s *TCPServer) Start() error {
	if s.Running.Load() {
		s.Logger.Error("server already running")
		return fmt.Errorf("server %s already running", s.Name)
	}

	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		s.Logger.Error("server failed to start", shared.Field{Key: "error", Value: err})
		return fmt.Errorf("server %s failed to start: %w", s.Name, err)
	}

	s.Listener = ln
	s.Running.Store(true)

	s.Logger.Info("server started", shared.Field{Key: "addr", Value: s.Addr})
	go s.AcceptLoop()

	return nil
}

func (s *TCPServer) Stop() {
	if !s.Running.Load() {
		s.Logger.Info("server not running")
		return
	}

	s.Running.Store(false)
	if s.Listener != nil {
		_ = s.Listener.Close()
	}

	s.Sessions.Range(func(key uint32, session TCPServerSession) bool {
		if closer, ok := any(session).(interface{ Close() error }); ok {
			_ = closer.Close()
			return true
		}

		return false
	})

	s.Logger.Info("server stopped")
}

func (s *TCPServer) AddSession(id uint32, session TCPServerSession) {
	s.Sessions.Store(id, session)
}

func (s *TCPServer) RemoveSession(id uint32) {
	s.Sessions.Delete(id)
}

func (s *TCPServer) GetSession(id uint32) (TCPServerSession, bool) {
	return s.Sessions.Get(id)
}

func (s *TCPServer) AcceptLoop() {
	for s.Running.Load() {
		conn, err := s.Listener.Accept()
		if err != nil {
			if !s.Running.Load() {
				return
			}

			s.Logger.Error("accept error", shared.Field{Key: "error", Value: err})
			continue
		}

		id := s.UidGenerator.Uid()
		session := s.NewSession(id, conn, s)
		s.AddSession(id, session)
		go session.Handle()
	}
}
