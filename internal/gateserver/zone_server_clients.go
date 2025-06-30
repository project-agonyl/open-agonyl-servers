package gateserver

import (
	"fmt"

	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/crypto"
)

type ZoneServerClients struct {
	servers map[byte]*ZoneServerClient
	crypto  crypto.Crypto
	players *Players
	logger  shared.Logger
}

func NewZoneServerClients(cfg *config.EnvVars, crypto crypto.Crypto, players *Players, logger shared.Logger) *ZoneServerClients {
	servers := make(map[byte]*ZoneServerClient)
	for _, zoneServer := range cfg.ZoneServers {
		server := NewZoneServerClient(
			byte(zoneServer.ID),
			cfg.ServerId,
			zoneServer.IP,
			uint32(zoneServer.Port),
			logger,
			players,
			crypto,
		)
		servers[byte(zoneServer.ID)] = server
	}

	return &ZoneServerClients{
		servers: servers,
		crypto:  crypto,
		players: players,
		logger:  logger,
	}
}

func (zs *ZoneServerClients) Start() {
	for _, server := range zs.servers {
		go server.Start()
	}
}

func (zs *ZoneServerClients) Stop() {
	for _, server := range zs.servers {
		server.Stop()
	}
}

func (zs *ZoneServerClients) GetServer(id byte) (*ZoneServerClient, error) {
	server, exists := zs.servers[id]
	if !exists {
		return nil, fmt.Errorf("zone server with id %d not found", id)
	}

	return server, nil
}

func (zs *ZoneServerClients) Send(id byte, packet []byte) error {
	server, err := zs.GetServer(id)
	if err != nil {
		return err
	}

	return server.Send(packet)
}
