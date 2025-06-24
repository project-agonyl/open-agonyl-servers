package accountserver

import (
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	dbService db.DBService
	cfg       *config.EnvVars
}

func NewServer(cfg *config.EnvVars, db db.DBService, logger shared.Logger) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         cfg.IpAddress + ":" + cfg.Port,
			Name:         "account-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		dbService: db,
		cfg:       cfg,
	}
	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newAccountServerSession(id, conn)
		if accountSession, ok := session.(*accountServerSession); ok {
			accountSession.server = server
		}

		return session
	}
	return server
}
