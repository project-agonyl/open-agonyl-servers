package mainserver

type PlayerState byte

const (
	PlayerStateLogin PlayerState = iota
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
}

func NewPlayer(
	pcId uint32,
	account string,
	characterName string,
	clientIp string,
	currentMapId uint16,
	serverId byte,
	gateServerId byte,
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
	}
}
