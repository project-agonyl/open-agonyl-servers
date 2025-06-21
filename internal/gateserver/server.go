package gateserver

import (
	"fmt"
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/crypto"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	cfg               *config.EnvVars
	players           *Players
	zoneServerClients *ZoneServerClients
	loginServerClient *LoginServerClient
	crypto            crypto.Crypto
	db                db.DBService
}

func NewServer(
	logger shared.Logger,
	cfg *config.EnvVars,
	players *Players,
	zoneServerClients *ZoneServerClients,
	loginServerClient *LoginServerClient,
	crypto crypto.Crypto,
	db db.DBService,
) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Name:         "gate-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		cfg:               cfg,
		players:           players,
		zoneServerClients: zoneServerClients,
		loginServerClient: loginServerClient,
		crypto:            crypto,
		db:                db,
	}
	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newServerSession(id, conn)
		if serverSession, ok := session.(*serverSession); ok {
			serverSession.server = server
		}

		return session
	}
	return server
}
