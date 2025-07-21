package zoneserver

import (
	"fmt"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type PlayerState int

const (
	PlayerStateWorldLoginPending PlayerState = iota
	PlayerStateWorldLoginSuccess
	PlayerStateInGame
)

type Player struct {
	PcId              uint32
	Account           string
	CharacterName     string
	Class             byte
	Level             uint16
	Exp               uint32
	Location          Location
	Skills            []Skill
	PKCount           uint32
	RTime             uint32
	SocialInfo        SocialInfo
	Woonz             uint32
	Lore              uint32
	Stats             Stats
	Wear              []WearItem
	Inventory         []InventoryItem
	ActivePet         Pet
	PetInventory      []PetInventory
	GateServerSession *zoneServerSession
	Logger            shared.Logger
	Zone              *Zone
	State             PlayerState
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
		PcId:              pcId,
		Account:           account,
		CharacterName:     characterName,
		GateServerSession: gateServerSession,
		Logger:            logger,
		Zone:              zone,
		State:             PlayerStateWorldLoginPending,
	}
}

func (p *Player) HandleGateServerPacket(packet []byte) {
	if len(packet) < 12 {
		return
	}

	if !p.Zone.EnqueuePlayerPacket(packet) {
		p.Logger.Error(
			"Failed to enqueue player packet",
			shared.Field{Key: "playerId", Value: p.PcId},
			shared.Field{Key: "packet", Value: packet},
		)
	}
}

func (p *Player) HandleMainServerPacket(packet []byte) {
	if len(packet) < 9 {
		return
	}

	if !p.Zone.EnqueueMainServerPacket(packet) {
		p.Logger.Error(
			"Failed to enqueue main server packet",
			shared.Field{Key: "playerId", Value: p.PcId},
			shared.Field{Key: "packet", Value: packet},
		)
	}
}

func (p *Player) Send(packet []byte) error {
	if p.GateServerSession == nil {
		return fmt.Errorf("gate server session is nil")
	}

	return p.GateServerSession.Send(packet)
}

type Location struct {
	MapId uint16
	X     byte
	Y     byte
}

type Skill struct {
	Id    byte
	Level byte
}

type SocialInfo struct {
	KHRank byte
	KHId   uint32
	KHName string
	Nation byte
}

type Stats struct {
	// Base stats
	RemainingPoints uint16
	Strength        uint16
	Intelligence    uint16
	Dexterity       uint16
	Vitality        uint16
	Mana            uint16
	HPCapacity      uint16
	MPCapacity      uint16
	HP              uint16
	MP              uint16

	// Calculated stats
	HitAttack             uint16
	MagicAttack           uint16
	Defense               uint16
	FireAttack            uint16
	FireDefence           uint16
	IceAttack             uint16
	IceDefense            uint16
	LightAttack           uint16
	LightDefense          uint16
	MaxHp                 uint16
	MaxMp                 uint16
	AdditionalHitAttack   uint16
	AdditionalMagicAttack uint16
}

type WearItem struct {
	ItemCode       uint32
	ItemOption     uint32
	ItemUniqueCode uint32
	WearIndex      byte
}

type InventoryItem struct {
	ItemCode       uint32
	ItemOption     uint32
	ItemUniqueCode uint32
	Slot           byte
}

type Pet struct {
	PetCode       uint32
	PetHP         uint32
	PetOption     uint32
	PetUniqueCode uint32
}

type PetInventory struct {
	Pet  Pet
	Slot byte
}
