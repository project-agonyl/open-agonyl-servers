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
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/messages/protocol"
	"github.com/project-agonyl/open-agonyl-servers/internal/utils"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/db"
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
	zoneManager     *ZoneManager
	db              db.DBService
}

func NewMainServerClient(
	serverId byte,
	addr string,
	logger shared.Logger,
	players *Players,
	zoneManager *ZoneManager,
	db db.DBService,
) *MainServerClient {
	return &MainServerClient{
		serverId:    serverId,
		addr:        addr,
		logger:      logger,
		players:     players,
		isConnected: false,
		zoneManager: zoneManager,
		db:          db,
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
	if proto == protocol.M2SWorldLogin {
		if exists {
			c.logger.Error(
				"Player already logged in",
				shared.Field{Key: "pcId", Value: pcId},
				shared.Field{Key: "protocol", Value: proto},
			)
			return
		}

		msg, err := messages.ReadMsgM2SWorldLogin(packet)
		if err != nil {
			c.logger.Error(
				"Failed to read M2SWorldLogin message",
				shared.Field{Key: "error", Value: err},
				shared.Field{Key: "packet", Value: packet},
			)
			return
		}

		characterName := utils.ReadStringFromBytes(msg.CharacterName[:])
		characterData, err := c.db.GetCharacter(pcId, characterName)
		if err != nil {
			c.logger.Error(
				"Failed to get character",
				shared.Field{Key: "error", Value: err},
				shared.Field{Key: "characterName", Value: characterName},
				shared.Field{Key: "pcId", Value: pcId},
			)
			return
		}

		player := NewPlayer(
			pcId,
			characterData.Account,
			characterName,
			nil,
			c.logger,
			c.zoneManager.GetZone(msg.MapId),
		)
		player.Class = characterData.Class
		player.Level = characterData.Level
		player.Lore = characterData.Data.Lore
		player.Woonz = characterData.Data.Parole
		player.SocialInfo = SocialInfo{
			Nation: characterData.Data.SocialInfo.Nation,
		}
		player.Location = Location{
			MapId: characterData.Data.Location.MapCode,
			X:     characterData.Data.Location.Position.X,
			Y:     characterData.Data.Location.Position.Y,
		}
		player.Stats = Stats{
			RemainingPoints: characterData.Data.Stats.RemainingPoints,
			Strength:        characterData.Data.Stats.Strength,
			Intelligence:    characterData.Data.Stats.Intelligence,
			Dexterity:       characterData.Data.Stats.Dexterity,
			Vitality:        characterData.Data.Stats.Vitality,
			Mana:            characterData.Data.Stats.Mana,
			HPCapacity:      characterData.Data.Stats.HPCapacity,
			MPCapacity:      characterData.Data.Stats.MPCapacity,
			HP:              characterData.Data.Stats.HP,
			MP:              characterData.Data.Stats.MP,
		}
		player.Wear = make([]WearItem, len(characterData.Data.Wear))
		for i, wearItem := range characterData.Data.Wear {
			player.Wear[i] = WearItem{
				ItemCode:       wearItem.ItemCode,
				ItemOption:     wearItem.ItemOption,
				ItemUniqueCode: wearItem.ItemUniqueCode,
			}
		}
		player.Inventory = make([]InventoryItem, len(characterData.Data.Inventory))
		for i, invItem := range characterData.Data.Inventory {
			player.Inventory[i] = InventoryItem{
				ItemCode:       invItem.ItemCode,
				ItemOption:     invItem.ItemOption,
				ItemUniqueCode: invItem.ItemUniqueCode,
				Slot:           invItem.Slot,
			}
		}
		player.Skills = make([]Skill, len(characterData.Data.Skills))
		for i, skill := range characterData.Data.Skills {
			player.Skills[i] = Skill{
				Id:    skill.SkillID,
				Level: skill.Level,
			}
		}
		player.ActivePet = Pet{
			PetCode:       characterData.Data.ActivePet.PetCode,
			PetHP:         characterData.Data.ActivePet.PetHP,
			PetOption:     characterData.Data.ActivePet.PetOption,
			PetUniqueCode: characterData.Data.ActivePet.PetUniqueCode,
		}
		player.PetInventory = make([]PetInventory, len(characterData.Data.PetInventory))
		for i, petInv := range characterData.Data.PetInventory {
			player.PetInventory[i] = PetInventory{
				Pet: Pet{
					PetCode:       petInv.PetCode,
					PetHP:         petInv.PetHP,
					PetOption:     petInv.PetOption,
					PetUniqueCode: petInv.PetUniqueCode,
				},
				Slot: petInv.Slot,
			}
		}

		player.Zone.EnqueuePlayerLogin(pcId)
		player.State = PlayerStateWorldLoginSuccess
		c.players.Add(player)
		return
	}

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
