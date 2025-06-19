package loginserver

import (
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
}

func NewServer(addr string, logger shared.Logger, cacheService shared.CacheService, dbService db.DBService, broker *Broker) *Server {
	return &Server{
		TCPServer: network.TCPServer{
			Addr:         addr,
			Name:         "login-server",
			NewSession:   newLoginServerSession,
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
		},
		cacheService: cacheService,
		dbService:    dbService,
		broker:       broker,
	}
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
