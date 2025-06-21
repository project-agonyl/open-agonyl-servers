package gateserver

import (
	"fmt"
	"net"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type Player struct {
	Id          uint32
	Username    string
	conn        net.Conn
	currentZone byte
	sendChan    chan []byte
	done        chan struct{}
	wg          sync.WaitGroup
	logger      shared.Logger
}

func NewPlayer(id uint32, username string, conn net.Conn, logger shared.Logger) *Player {
	player := &Player{
		Id:          id,
		Username:    username,
		conn:        conn,
		currentZone: 255,
		sendChan:    make(chan []byte, 100),
		done:        make(chan struct{}),
		logger:      logger,
	}

	player.wg.Add(1)
	go player.sender()

	return player
}

func (s *Player) Send(data []byte) error {
	select {
	case s.sendChan <- data:
		return nil
	case <-s.done:
		return fmt.Errorf("session is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (p *Player) Close() {
	close(p.done)
	p.wg.Wait()
}

func (p *Player) GetCurrentZone() byte {
	return p.currentZone
}

func (p *Player) sender() {
	defer p.wg.Done()
	for {
		select {
		case data := <-p.sendChan:
			if data[8] == 0x01 && data[9] == 0xE1 {
				if len(data) > 0x0A {
					p.currentZone = data[0x0A]
				}

				continue
			}

			if _, err := p.conn.Write(data); err != nil {
				p.logger.Error("Failed to send packet to client",
					shared.Field{Key: "error", Value: err},
					shared.Field{Key: "playerId", Value: p.Id},
					shared.Field{Key: "currentZone", Value: p.currentZone},
					shared.Field{Key: "username", Value: p.Username},
				)
				return
			}

		case <-p.done:
			return
		}
	}
}
