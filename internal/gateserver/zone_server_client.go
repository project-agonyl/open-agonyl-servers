package gateserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/crypto"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
)

type ZoneServerClient struct {
	id              byte
	name            string
	ip              string
	port            uint32
	addr            string
	conn            net.Conn
	running         atomic.Bool
	shouldReconnect atomic.Bool
	sendChan        chan []byte
	done            chan struct{}
	wg              sync.WaitGroup
	logger          shared.Logger
	reconnectDelay  time.Duration
	players         *Players
	crypto          crypto.Crypto
	isConnected     bool
}

func NewZoneServerClient(id byte, ip string, port uint32, logger shared.Logger, players *Players, crypto crypto.Crypto) *ZoneServerClient {
	return &ZoneServerClient{
		id:             id,
		ip:             ip,
		port:           port,
		addr:           fmt.Sprintf("%s:%d", ip, port),
		name:           fmt.Sprintf("zone server %d", id),
		logger:         logger,
		reconnectDelay: 10 * time.Second,
		players:        players,
		crypto:         crypto,
		isConnected:    false,
	}
}

func (c *ZoneServerClient) Start() {
	c.running.Store(true)
	c.shouldReconnect.Store(true)
	c.logger.Info(
		fmt.Sprintf("Starting %s client", c.name),
		shared.Field{Key: "addr", Value: c.addr},
		shared.Field{Key: "name", Value: c.name},
		shared.Field{Key: "serverId", Value: c.id},
	)
	for c.running.Load() {
		if err := c.connect(); err != nil {
			if !c.shouldReconnect.Load() {
				break
			}

			time.Sleep(c.reconnectDelay)
			continue
		}

		c.handleConnection()
		if !c.shouldReconnect.Load() {
			break
		}

		time.Sleep(c.reconnectDelay)
	}
}

func (c *ZoneServerClient) Stop() {
	c.running.Store(false)
	c.shouldReconnect.Store(false)
	c.logger.Info(
		fmt.Sprintf("Stopping %s client", c.name),
		shared.Field{Key: "addr", Value: c.addr},
		shared.Field{Key: "name", Value: c.name},
		shared.Field{Key: "serverId", Value: c.id},
	)
}

func (c *ZoneServerClient) Send(packet []byte) error {
	if !c.isConnected {
		return fmt.Errorf("client is not connected")
	}

	select {
	case c.sendChan <- packet:
		return nil
	case <-c.done:
		return fmt.Errorf("client is closing")
	default:
		return fmt.Errorf("send channel is full")
	}
}

func (c *ZoneServerClient) connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.sendChan = make(chan []byte, 100)
	c.done = make(chan struct{})
	c.isConnected = true
	go c.logger.Info(
		fmt.Sprintf("Connected to %s", c.name),
		shared.Field{Key: "addr", Value: c.addr},
		shared.Field{Key: "name", Value: c.name},
		shared.Field{Key: "serverId", Value: c.id},
	)

	return nil
}

func (c *ZoneServerClient) handleConnection() {
	defer func() {
		close(c.done)
		c.wg.Wait()
		_ = c.conn.Close()
		c.conn = nil
		c.isConnected = false
	}()
	c.wg.Add(1)
	go c.sender()
	msg := messages.NewMsgGate2ZsConnect(c.id)
	err := c.Send(msg.GetBytes())
	if err != nil {
		c.logger.Error(
			fmt.Sprintf("Failed to send connect packet to %s", c.name),
			shared.Field{Key: "addr", Value: c.addr},
			shared.Field{Key: "name", Value: c.name},
			shared.Field{Key: "serverId", Value: c.id},
			shared.Field{Key: "error", Value: err},
		)
		return
	}

	for {
		var buf bytes.Buffer
		if _, err := io.CopyN(&buf, c.conn, 4); err != nil {
			break
		}

		reader := io.MultiReader(&buf, c.conn)
		dataLength := binary.LittleEndian.Uint32(buf.Bytes())
		if dataLength == 0 {
			continue
		}

		if dataLength > 16*1024*1024 {
			break
		}

		packet := make([]byte, dataLength)
		if _, err := io.ReadFull(reader, packet); err != nil {
			break
		}

		c.processPacket(packet)
	}
}

func (c *ZoneServerClient) processPacket(packet []byte) {
	pcId := binary.LittleEndian.Uint32(packet[4:])
	player, exists := c.players.Get(pcId)
	if !exists {
		return
	}

	go func(p []byte) {
		if p[8] != 0x01 && p[9] != 0xE1 {
			c.crypto.Encrypt(p)
		}

		_ = player.Send(p)
	}(packet)
}

func (c *ZoneServerClient) sender() {
	defer c.wg.Done()
	for {
		select {
		case data := <-c.sendChan:
			if _, err := c.conn.Write(data); err != nil {
				c.logger.Error(
					fmt.Sprintf("Failed to send packet to %s", c.name),
					shared.Field{Key: "addr", Value: c.addr},
					shared.Field{Key: "name", Value: c.name},
					shared.Field{Key: "error", Value: err},
					shared.Field{Key: "serverId", Value: c.id})
				return
			}

		case <-c.done:
			return
		}
	}
}
