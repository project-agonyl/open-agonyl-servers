package mainserver

type PlayerState byte

const (
	PlayerStateLogin PlayerState = iota
	PlayerStateWorld PlayerState = 1
)

type Player struct {
	pcId            uint32
	account         string
	characterName   string
	clientIp        string
	currentMapId    uint16
	state           PlayerState
	currentServerId byte
	gateServerId    byte
	zone            *Zone
}

func NewPlayer(
	pcId uint32,
	account string,
	characterName string,
	clientIp string,
	currentMapId uint16,
	serverId byte,
	gateServerId byte,
	zone *Zone,
) *Player {
	return &Player{
		pcId:            pcId,
		account:         account,
		characterName:   characterName,
		clientIp:        clientIp,
		currentMapId:    currentMapId,
		state:           PlayerStateLogin,
		currentServerId: serverId,
		gateServerId:    gateServerId,
		zone:            zone,
	}
}
