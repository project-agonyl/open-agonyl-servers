package zoneserver

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

func (p *Players) GetByCharacterName(characterName string) (player *Player, exists bool) {
	p.players.Range(func(key uint32, value *Player) bool {
		if value.characterName == characterName {
			player = value
			exists = true
			return false
		}

		return true
	})

	return player, exists
}

func (p *Players) GetByAccount(account string) (player *Player, exists bool) {
	p.players.Range(func(key uint32, value *Player) bool {
		if value.account == account {
			player = value
			exists = true
			return false
		}

		return true
	})

	return player, exists
}
