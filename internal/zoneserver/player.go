package zoneserver

import (
	"encoding/binary"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
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
}

func NewPlayer(
	pcId uint32,
	account string,
	characterName string,
	gateServerSession *zoneServerSession,
	logger shared.Logger,
) *Player {
	return &Player{
		pcId:              pcId,
		account:           account,
		characterName:     characterName,
		gateServerSession: gateServerSession,
		logger:            logger,
	}
}

func (p *Player) HandleGateServerPacket(packet []byte) {
	if len(packet) < 12 {
		return
	}

	proto := binary.LittleEndian.Uint16(packet[10:])
	switch proto {
	default:
		p.logger.Error(
			"Unhandled gate server packet",
			shared.Field{Key: "protocol", Value: proto},
		)
	}
}

func (p *Player) HandleMainServerPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	proto := binary.LittleEndian.Uint16(packet)
	switch proto {
	case protocol.M2SError:
		msg, err := messages.ReadMsgM2SError(packet)
		if err != nil {
			return
		}

		message := utils.ReadStringFromBytes(msg.Msg[:])
		_ = p.gateServerSession.SendErrorMsg(p.pcId, msg.Code, message)
	default:
		p.logger.Error(
			"Unhandled main server packet",
			shared.Field{Key: "protocol", Value: proto},
		)
	}
}
