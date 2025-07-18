package zoneserver

import (
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type PlayerState int

const (
	PlayerStateWorldLoginPending PlayerState = iota
	PlayerStateWorldLoginSuccess
)

type Player struct {
	pcId              uint32
	account           string
	characterName     string
	gateServerSession *zoneServerSession
	logger            shared.Logger
	zone              *Zone
}

func NewPlayer(
	pcId uint32,
	account string,
	characterName string,
	gateServerSession *zoneServerSession,
	logger shared.Logger,
	zone *Zone,
) *Player {
	return &Player{
		pcId:              pcId,
		account:           account,
		characterName:     characterName,
		gateServerSession: gateServerSession,
		logger:            logger,
		zone:              zone,
	}
}

func (p *Player) HandleGateServerPacket(packet []byte) {
	if len(packet) < 12 {
		return
	}

	if !p.zone.EnqueuePlayerPacket(packet) {
		p.logger.Error(
			"Failed to enqueue player packet",
			shared.Field{Key: "playerId", Value: p.pcId},
			shared.Field{Key: "packet", Value: packet},
		)
	}
}

func (p *Player) HandleMainServerPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	if !p.zone.EnqueueMainServerPacket(packet) {
		p.logger.Error(
			"Failed to enqueue main server packet",
			shared.Field{Key: "playerId", Value: p.pcId},
			shared.Field{Key: "packet", Value: packet},
		)
	}
}
