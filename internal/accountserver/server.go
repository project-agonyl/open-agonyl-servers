package accountserver

import (
	"errors"
	"net"

	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/data"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Server struct {
	network.TCPServer
	dbService        db.DBService
	cfg              *config.EnvVars
	items            map[uint32]data.Item
	players          *Players
	mainServerClient *MainServerClient
}

func NewServer(cfg *config.EnvVars, db db.DBService, logger shared.Logger, players *Players, mainServerClient *MainServerClient) *Server {
	server := &Server{
		TCPServer: network.TCPServer{
			Addr:         cfg.IpAddress + ":" + cfg.Port,
			Name:         "account-server",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		dbService:        db,
		cfg:              cfg,
		items:            make(map[uint32]data.Item),
		players:          players,
		mainServerClient: mainServerClient,
	}
	server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newAccountServerSession(id, conn)
		if accountSession, ok := session.(*accountServerSession); ok {
			accountSession.server = server
		}

		return session
	}
	err := server.loadItems()
	if err != nil {
		logger.Error("Failed to load items", shared.Field{Key: "error", Value: err})
	}

	return server
}

func (s *Server) GetItem(itemCode uint32) (*data.Item, bool) {
	item, ok := s.items[itemCode]
	if !ok {
		return nil, false
	}

	return &item, true
}

func (s *Server) loadItems() error {
	if s.cfg.ZoneDataItemPath == "" {
		return errors.New("ZoneDataItemPath is not set")
	}

	items, err := data.LoadIT0Items(s.cfg.ZoneDataItemPath+"/0", s.cfg.ZoneDataItemPath+"/0ex")
	if err == nil {
		for _, item := range items {
			s.items[item.ItemCode] = item
		}
	}

	items, err = data.LoadIT1Items(s.cfg.ZoneDataItemPath + "/1")
	if err == nil {
		for _, item := range items {
			s.items[item.ItemCode] = item
		}
	}

	return err
}
