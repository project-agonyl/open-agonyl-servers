package mainserver

import (
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/mainserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/mainserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	cfg          *config.EnvVars
	dbService    db.DBService
	players      *Players
	mapZones     *shared.SafeMap[uint16, *Zone]
	zoneSessions *shared.SafeMap[byte, *Zone]
}

func NewServer(cfg *config.EnvVars, db db.DBService, logger shared.Logger, players *Players) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         cfg.IpAddress + ":" + cfg.Port,
			Name:         "main-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		cfg:          cfg,
		dbService:    db,
		players:      players,
		mapZones:     shared.NewSafeMap[uint16, *Zone](),
		zoneSessions: shared.NewSafeMap[byte, *Zone](),
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

func (s *Server) IsZoneRegistered(zoneId byte) bool {
	return s.zoneSessions.Has(zoneId)
}
