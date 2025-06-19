package loginserver

import (
	"slices"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/helpers"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/network"
)

type Broker struct {
	network.TCPServer
	cacheService shared.CacheService
}

func NewBroker(addr string, logger shared.Logger, cacheService shared.CacheService) *Broker {
	return &Broker{
		TCPServer: network.TCPServer{
			Addr:         addr,
			Name:         "login-server-broker",
			NewSession:   newBrokerSession,
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
		},
		cacheService: cacheService,
	}
}

func (s *Broker) GetGateServerList() map[byte]string {
	addedPorts := []uint32{}
	result := map[byte]string{}
	s.Sessions.Range(func(key uint32, value network.TCPServerSession) bool {
		brokerSession := value.(*brokerSession)
		if !slices.Contains(addedPorts, brokerSession.port) {
			addedPorts = append(addedPorts, brokerSession.port)
			result[brokerSession.serverId] = brokerSession.serverName
		}

		return true
	})

	return result
}

func (s *Broker) GetGateServerCount() int {
	return s.Sessions.Len()
}

func (s *Broker) AddLoggedInUser(username string, id uint32) {
	helpers.AddLoggedInUser(s.cacheService, username, id)
}

func (s *Broker) RemoveLoggedInUser(username string) {
	helpers.RemoveLoggedInUser(s.cacheService, username)
}

func (s *Broker) IsLoggedIn(username string) bool {
	return helpers.IsLoggedIn(s.cacheService, username)
}
