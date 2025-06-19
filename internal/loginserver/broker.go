package loginserver

import (
	"errors"
	"net"
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
	broker := &Broker{
		TCPServer: network.TCPServer{
			Addr:         addr,
			Name:         "login-broker",
			UidGenerator: shared.NewUidGenerator(0),
			Logger:       logger,
			Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
		},
		cacheService: cacheService,
	}
	broker.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
		session := newBrokerSession(id, conn)
		if brokerSession, ok := session.(*brokerSession); ok {
			brokerSession.server = broker
		}

		return session
	}

	return broker
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

func (s *Broker) GetGateServerInfoByServerId(serverId byte) (*GateInfo, error) {
	gateInfo := GateInfo{}
	found := false
	s.Sessions.Range(func(key uint32, value network.TCPServerSession) bool {
		brokerSession := value.(*brokerSession)
		if brokerSession.serverId == serverId {
			gateInfo.Id = brokerSession.id
			gateInfo.IpAddress = brokerSession.ipAddress
			gateInfo.Port = brokerSession.port
			found = true
			return false
		}

		return true
	})

	if !found {
		return nil, errors.New("gate server not found")
	}

	return &gateInfo, nil
}

func (s *Broker) SendMsgToGateServer(id uint32, msg []byte) error {
	session, ok := s.Sessions.Get(id)
	if !ok {
		return errors.New("gate server not found")
	}

	return session.Send(msg)
}

type GateInfo struct {
	Id        uint32
	IpAddress string
	Port      uint32
}
