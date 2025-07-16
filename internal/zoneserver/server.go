package zoneserver

import (
	"fmt"
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/db"
)

type Server struct {
	network.TCPServer
	Logger                shared.Logger
	db                    db.DBService
	mainServerClient      *MainServerClient
	serialNumberGenerator shared.SerialNumberGenerator
	players               *Players
}

func NewServer(
	cfg *config.EnvVars,
	db db.DBService,
	logger shared.Logger,
	mainServerClient *MainServerClient,
	serialNumberGenerator shared.SerialNumberGenerator,
	players *Players,
) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         cfg.IpAddress + ":" + cfg.Port,
			Name:         fmt.Sprintf("zone-server-%d", cfg.ServerId),
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		db:                    db,
		mainServerClient:      mainServerClient,
		serialNumberGenerator: serialNumberGenerator,
		players:               players,
	}

	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newZoneServerSession(id, conn)
		if zoneSession, ok := session.(*zoneServerSession); ok {
			zoneSession.server = server
		}

		return session
	}

	return server
}
