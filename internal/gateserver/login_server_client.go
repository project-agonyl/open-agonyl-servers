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

	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
)

type LoginServerClient struct {
	id               byte
	name             string
	serverName       string
	serverIpAddress  string
	serverPort       uint32
	ipAddress        string
	port             uint32
	addr             string
	conn             net.Conn
	running          atomic.Bool
	shouldReconnect  atomic.Bool
	sendChan         chan []byte
	done             chan struct{}
	wg               sync.WaitGroup
	logger           shared.Logger
	reconnectDelay   time.Duration
	loggedInAccounts *shared.SafeMap[uint32, string]
	isConnected      bool
}

func NewLoginServerClient(cfg *config.EnvVars, logger shared.Logger) *LoginServerClient {
	return &LoginServerClient{
		id:               cfg.ServerId,
		ipAddress:        cfg.LoginServerIpAddress,
		port:             uint32(cfg.GetLoginServerPort()),
		addr:             fmt.Sprintf("%s:%d", cfg.LoginServerIpAddress, cfg.GetLoginServerPort()),
		name:             "login server",
		serverName:       cfg.ServerName,
		serverPort:       uint32(cfg.GetServerPort()),
		serverIpAddress:  cfg.IpAddress,
		logger:           logger,
		reconnectDelay:   10 * time.Second,
		loggedInAccounts: shared.NewSafeMap[uint32, string](),
		isConnected:      false,
	}
}

func (c *LoginServerClient) Start() {
	c.running.Store(true)
	c.shouldReconnect.Store(true)
	c.logger.Info(
		fmt.Sprintf("Starting %s client", c.name),
		shared.Field{Key: "addr", Value: c.addr},
		shared.Field{Key: "name", Value: c.name},
		shared.Field{Key: "clientId", Value: c.id},
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

func (c *LoginServerClient) Stop() {
	c.running.Store(false)
	c.shouldReconnect.Store(false)
	c.logger.Info(
		fmt.Sprintf("Stopping %s client", c.name),
		shared.Field{Key: "addr", Value: c.addr},
		shared.Field{Key: "name", Value: c.name},
		shared.Field{Key: "clientId", Value: c.id},
	)
}

func (c *LoginServerClient) Send(packet []byte) error {
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

func (c *LoginServerClient) IsLoggedIn(pcId uint32) bool {
	_, ok := c.loggedInAccounts.Get(pcId)
	return ok
}

func (c *LoginServerClient) GetLoggedInAccount(pcId uint32) string {
	account, ok := c.loggedInAccounts.Get(pcId)
	if !ok {
		return ""
	}

	return account
}

func (c *LoginServerClient) RemoveLoggedInAccount(pcId uint32) {
	c.loggedInAccounts.Delete(pcId)
}

func (c *LoginServerClient) PopLoggedInAccount(pcId uint32) (string, bool) {
	account, ok := c.loggedInAccounts.Get(pcId)
	if !ok {
		return "", false
	}

	c.loggedInAccounts.Delete(pcId)
	return account, true
}

func (c *LoginServerClient) connect() error {
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
		shared.Field{Key: "clientId", Value: c.id},
	)

	return nil
}

func (c *LoginServerClient) handleConnection() {
	defer func() {
		close(c.done)
		c.wg.Wait()
		_ = c.conn.Close()
		c.conn = nil
		c.isConnected = false
	}()
	c.wg.Add(1)
	go c.sender()
	msg := messages.NewMsgGate2LsConnect(c.id, c.id, c.serverIpAddress, c.serverPort, c.serverName)
	err := c.Send(msg.GetBytes())
	if err != nil {
		c.logger.Error(
			fmt.Sprintf("Failed to send connect packet to %s", c.name),
			shared.Field{Key: "addr", Value: c.addr},
			shared.Field{Key: "name", Value: c.name},
			shared.Field{Key: "clientId", Value: c.id},
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

func (c *LoginServerClient) processPacket(packet []byte) {
	ctrl := packet[8]
	cmd := packet[9]
	switch ctrl {
	case 0x01:
		switch cmd {
		case 0xE1:
			c.handleLogin(packet)
		case 0xE3:
			c.handleLogout(packet)
		}
	}
}

func (c *LoginServerClient) handleLogin(packet []byte) {
	msg, err := messages.ReadMsgLs2GateLogin(packet)
	if err != nil {
		return
	}

	c.loggedInAccounts.Set(msg.PcId, string(msg.Account[:]))
}

func (c *LoginServerClient) handleLogout(packet []byte) {
	// TODO: Implement
}

func (c *LoginServerClient) sender() {
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
					shared.Field{Key: "clientId", Value: c.id})
				return
			}

		case <-c.done:
			return
		}
	}
}
