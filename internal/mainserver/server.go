package mainserver

import (
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/mainserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	cfg *config.EnvVars
}

func NewServer(cfg *config.EnvVars, logger shared.Logger) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         cfg.IpAddress + ":" + cfg.Port,
			Name:         "main-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		cfg: cfg,
	}
	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newMainServerSession(id, conn)
		if mainSession, ok := session.(*mainServerSession); ok {
			mainSession.server = server
		}

		return session
	}

	return server
}
