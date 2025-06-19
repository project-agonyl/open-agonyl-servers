package loginserver

import (
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/helpers"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	dbService    db.DBService
	cacheService shared.CacheService
	broker       *Broker
	cfg          *config.EnvVars
}

func NewServer(
	addr string,
	logger shared.Logger,
	cacheService shared.CacheService,
	dbService db.DBService,
	broker *Broker,
	cfg *config.EnvVars,
) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         addr,
			Name:         "login-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		cacheService: cacheService,
		dbService:    dbService,
		broker:       broker,
		cfg:          cfg,
	}
	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newLoginServerSession(id, conn)
		if loginSession, ok := session.(*loginServerSession); ok {
			loginSession.server = server
		}

		return session
	}
	return server
}

func (s *Server) AddLoggedInUser(username string, id uint32) {
	helpers.AddLoggedInUser(s.cacheService, username, id)
}

func (s *Server) RemoveLoggedInUser(username string) {
	helpers.RemoveLoggedInUser(s.cacheService, username)
}

func (s *Server) IsLoggedIn(username string) bool {
	return helpers.IsLoggedIn(s.cacheService, username)
}
