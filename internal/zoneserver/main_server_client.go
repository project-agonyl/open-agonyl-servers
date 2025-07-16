package zoneserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

type MainServerClient struct {
	serverId        byte
	addr            string
	conn            net.Conn
	running         atomic.Bool
	shouldReconnect atomic.Bool
	sendChan        chan []byte
	done            chan struct{}
	wg              sync.WaitGroup
	logger          shared.Logger
	reconnectDelay  time.Duration
	isConnected     bool
	players         *Players
}

func NewMainServerClient(serverId byte, addr string, logger shared.Logger, players *Players) *MainServerClient {
	return &MainServerClient{
		serverId:    serverId,
		addr:        addr,
		logger:      logger,
		players:     players,
		isConnected: false,
	}
}

func (c *MainServerClient) Start() {
	c.running.Store(true)
	c.shouldReconnect.Store(true)
	c.logger.Info(
		"Starting main server client",
		shared.Field{Key: "addr", Value: c.addr},
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

func (c *MainServerClient) Send(packet []byte) error {
	if !c.isConnected {
		return fmt.Errorf("main server client is not connected")
	}

	select {
	case c.sendChan <- packet:
		return nil
	case <-c.done:
		return fmt.Errorf("main server client is closing")
	default:
		return fmt.Errorf("main server client send channel is full")
	}
}

func (c *MainServerClient) Stop() {
	c.running.Store(false)
	c.shouldReconnect.Store(false)
	c.logger.Info(
		"Stopping main server client",
		shared.Field{Key: "addr", Value: c.addr},
	)
}

func (c *MainServerClient) connect() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.sendChan = make(chan []byte, 100)
	c.done = make(chan struct{})
	c.isConnected = true
	go c.logger.Info(
		"Connected to main server",
		shared.Field{Key: "addr", Value: c.addr},
	)

	return nil
}

func (c *MainServerClient) handleConnection() {
	defer func() {
		close(c.done)
		c.wg.Wait()
		_ = c.conn.Close()
		c.conn = nil
		c.isConnected = false
		c.logger.Info(
			"Disconnected from main server",
			shared.Field{Key: "addr", Value: c.addr},
		)
	}()
	c.wg.Add(1)
	go c.sender()
	_ = c.Send([]byte{0x01, 0xA0, 0x00, 0x00, c.serverId})
	buffer := make([]byte, 1024*16)
	dynamicBuffer := bytes.NewBuffer(nil)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			break
		}

		dynamicBuffer.Write(buffer[:n])
		for dynamicBuffer.Len() >= 4 {
			dataLength := int(binary.LittleEndian.Uint16(dynamicBuffer.Bytes()[2:]))
			if dataLength > dynamicBuffer.Len() || dataLength == 0 {
				break
			}

			currentPacket := dynamicBuffer.Next(dataLength)
			go c.processPacket(currentPacket)
		}
	}
}

func (c *MainServerClient) sender() {
	defer c.wg.Done()
	for {
		select {
		case data := <-c.sendChan:
			if _, err := c.conn.Write(data); err != nil {
				c.logger.Error(
					"Failed to send packet to main server",
					shared.Field{Key: "addr", Value: c.addr})
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *MainServerClient) processPacket(packet []byte) {
	proto := binary.LittleEndian.Uint16(packet)
	pcId := binary.LittleEndian.Uint32(packet[4:])
	player, exists := c.players.Get(pcId)
	if !exists {
		c.logger.Error(
			"Could not find player",
			shared.Field{Key: "pcId", Value: pcId},
			shared.Field{Key: "protocol", Value: proto},
		)
		return
	}

	player.HandleMainServerPacket(packet)
}
