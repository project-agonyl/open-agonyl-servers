package accountserver

import "github.com/project-agonyl/open-agonyl-servers/internal/shared"

type Players struct {
	players *shared.SafeMap[uint32, *Player]
}

func NewPlayers() *Players {
	return &Players{
		players: shared.NewSafeMap[uint32, *Player](),
	}
}

func (p *Players) Add(player *Player) {
	p.players.Set(player.pcId, player)
}

func (p *Players) Remove(id uint32) {
	p.players.Delete(id)
}

func (p *Players) Get(id uint32) (*Player, bool) {
	return p.players.Get(id)
}

func (p *Players) HasPlayer(id uint32) bool {
	_, exists := p.players.Get(id)
	return exists
}
