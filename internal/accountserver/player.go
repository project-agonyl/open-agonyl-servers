package accountserver

type Player struct {
	pcId     uint32
	account  string
	clientIp string
}

func NewPlayer(pcId uint32, account string, clientIp string) *Player {
	return &Player{
		pcId:     pcId,
		account:  account,
		clientIp: clientIp,
	}
}
