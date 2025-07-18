package zoneserver

import (
	"errors"
	"strconv"
	"sync"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/data"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/db"
	"github.com/redis/go-redis/v9"
)

type ZoneManager struct {
	cfg                   *config.EnvVars
	db                    db.DBService
	logger                shared.Logger
	zones                 map[uint16]*Zone
	npcsData              *shared.SafeMap[uint16, *data.NPCData]
	itemsData             map[uint32]*data.Item
	serialNumberGenerator shared.SerialNumberGenerator
	players               *Players
	zoneWg                sync.WaitGroup
}

func NewZoneManager(
	cfg *config.EnvVars,
	db db.DBService,
	logger shared.Logger,
	redis *redis.Client,
	serialNumberGenerator shared.SerialNumberGenerator,
	players *Players,
) *ZoneManager {
	return &ZoneManager{
		cfg:                   cfg,
		db:                    db,
		logger:                logger,
		zones:                 make(map[uint16]*Zone, len(cfg.MapIDs)),
		npcsData:              shared.NewSafeMap[uint16, *data.NPCData](),
		itemsData:             make(map[uint32]*data.Item),
		serialNumberGenerator: serialNumberGenerator,
		players:               players,
	}
}

func (m *ZoneManager) Start() error {
	m.logger.Info("Loading IT0 data...")
	it0, err := data.LoadIT0Items(m.cfg.ZoneDataItemPath+"/0", m.cfg.ZoneDataItemPath+"/0ex")
	if err != nil {
		m.logger.Error(
			"Error loading IT0 data",
			shared.Field{Key: "error", Value: err},
		)
		return err
	}

	m.logger.Info("Loaded IT0 items", shared.Field{Key: "count", Value: len(it0)})
	for _, item := range it0 {
		m.itemsData[item.ItemCode] = &item
	}

	m.logger.Info("Loading IT1 data...")
	it1, err := data.LoadIT1Items(m.cfg.ZoneDataItemPath + "/1")
	if err != nil {
		m.logger.Error(
			"Error loading IT1 data",
			shared.Field{Key: "error", Value: err},
		)
		return err
	}

	m.logger.Info("Loaded IT1 items", shared.Field{Key: "count", Value: len(it1)})
	for _, item := range it1 {
		m.itemsData[item.ItemCode] = &item
	}

	m.logger.Info("Loading IT2 data...")
	it2, err := data.LoadIT2Items(m.cfg.ZoneDataItemPath + "/2")
	if err != nil {
		m.logger.Error(
			"Error loading IT2 data",
			shared.Field{Key: "error", Value: err},
		)
		return err
	}

	for _, item := range it2 {
		m.itemsData[item.ItemCode] = &item
	}

	m.logger.Info("Loading IT3 data...")
	it3, err := data.LoadIT3Items(m.cfg.ZoneDataItemPath + "/3")
	if err != nil {
		m.logger.Error(
			"Error loading IT3 data",
			shared.Field{Key: "error", Value: err},
		)
		return err
	}

	m.logger.Info("Loaded IT3 items", shared.Field{Key: "count", Value: len(it3)})
	for _, item := range it3 {
		m.itemsData[item.ItemCode] = &item
	}

	m.logger.Info("Loading zones...")
	for _, mapId := range m.cfg.MapIDs {
		zone, err := NewZone(m.cfg, m.db, m.logger, mapId, m.players, m)
		if err != nil {
			m.logger.Error(
				"Error creating zone",
				shared.Field{Key: "mapId", Value: mapId},
				shared.Field{Key: "error", Value: err},
			)
			return err
		}

		m.zoneWg.Add(1)
		go func() {
			defer m.zoneWg.Done()
			err := zone.Start()
			if err != nil {
				m.logger.Error(
					"Error starting zone",
					shared.Field{Key: "mapId", Value: mapId},
					shared.Field{Key: "error", Value: err},
				)
				panic(err)
			}
		}()

		m.zones[mapId] = zone
	}

	m.logger.Info("Loaded zones", shared.Field{Key: "count", Value: len(m.cfg.MapIDs)})
	m.zoneWg.Wait()
	return nil
}

func (m *ZoneManager) Stop() {
	for _, zone := range m.zones {
		zone.Stop()
	}
}

func (m *ZoneManager) GetZone(mapId uint16) *Zone {
	return m.zones[mapId]
}

func (m *ZoneManager) GetZoneByPlayer(player *Player) *Zone {
	return m.zones[player.zone.mapId]
}

func (m *ZoneManager) GetZoneByPlayerId(playerId uint32) *Zone {
	player, exists := m.players.Get(playerId)
	if !exists {
		return nil
	}

	return m.zones[player.zone.mapId]
}

func (m *ZoneManager) GetNPCData(npcId uint16) (*data.NPCData, error) {
	npcData, exists := m.npcsData.Get(npcId)
	if !exists {
		npcData, err := data.LoadNPCData(m.cfg.ZoneDataNPCPath + "/" + strconv.Itoa(int(npcId)))
		if err != nil {
			return nil, err
		}

		m.npcsData.Set(npcId, npcData)
		return npcData, nil
	}

	return npcData, nil
}

func (m *ZoneManager) GetItemData(itemCode uint32) (*data.Item, error) {
	itemData, exists := m.itemsData[itemCode]
	if !exists {
		return nil, errors.New("item data not found")
	}

	return itemData, nil
}

func (z *ZoneManager) EnqueuePlayerPacket(mapId uint16, packet []byte) bool {
	zone, exists := z.zones[mapId]
	if !exists {
		return false
	}

	return zone.EnqueuePlayerPacket(packet)
}

func (z *ZoneManager) EnqueueMainServerPacket(mapId uint16, packet []byte) bool {
	zone, exists := z.zones[mapId]
	if !exists {
		return false
	}

	return zone.EnqueueMainServerPacket(packet)
}
