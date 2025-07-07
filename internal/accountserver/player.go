package accountserver

import (
	"errors"
	"sync"
)

type Player struct {
	pcId                  uint32
	account               string
	clientIp              string
	session               *accountServerSession
	selectedCharacterName string
	selectedCharacterMu   sync.RWMutex
}

func NewPlayer(pcId uint32, account string, clientIp string, session *accountServerSession) *Player {
	return &Player{
		pcId:     pcId,
		account:  account,
		clientIp: clientIp,
		session:  session,
	}
}

func (p *Player) GetSelectedCharacterName() string {
	p.selectedCharacterMu.RLock()
	defer p.selectedCharacterMu.RUnlock()
	return p.selectedCharacterName
}

func (p *Player) SetSelectedCharacterName(characterName string) {
	p.selectedCharacterMu.Lock()
	defer p.selectedCharacterMu.Unlock()
	p.selectedCharacterName = characterName
}

func (p *Player) Send(data []byte) error {
	if p.session == nil {
		return errors.New("player session not found")
	}

	return p.session.Send(data)
}
