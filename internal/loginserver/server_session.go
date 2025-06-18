package loginserver

import (
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type loginServerSession struct {
	server *LoginServer
	conn   net.Conn
	id     uint32
}

func newLoginServerSession(id uint32, conn net.Conn, server interface{}) network.TCPServerSession {
	return &loginServerSession{
		id:     id,
		conn:   conn,
		server: server.(*LoginServer),
	}
}

func (s *loginServerSession) ID() uint32 {
	return s.id
}

func (s *loginServerSession) Handle() {
	panic("not implemented")
}

func (s *loginServerSession) Send(data []byte) error {
	panic("not implemented")
}

func (s *loginServerSession) Close() error {
	panic("not implemented")
}
