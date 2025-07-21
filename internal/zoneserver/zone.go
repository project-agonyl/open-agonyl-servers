package zoneserver

import (
	"strconv"
	"sync/atomic"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/data"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/db"
)

type Zone struct {
	serverId              byte
	mapId                 uint16
	players               *Players
	currentPlayers        []uint32
	logger                shared.Logger
	cfg                   *config.EnvVars
	db                    db.DBService
	zoneManager           *ZoneManager
	mapData               *data.MapData
	isRunning             atomic.Bool
	playerPacketQueue     *shared.SafeQueue[[]byte]
	mainServerPacketQueue *shared.SafeQueue[[]byte]
	playerLoginQueue      *shared.SafeQueue[uint32]
}

func NewZone(
	cfg *config.EnvVars,
	db db.DBService,
	logger shared.Logger,
	mapId uint16,
	players *Players,
	zoneManager *ZoneManager,
) (*Zone, error) {
	mapData, err := data.LoadMapData(cfg.ZoneDataMapPath + "/" + strconv.Itoa(int(mapId)))
	if err != nil {
		return nil, err
	}

	return &Zone{
		serverId:              cfg.ServerId,
		mapId:                 mapId,
		players:               players,
		currentPlayers:        make([]uint32, 0),
		logger:                logger,
		cfg:                   cfg,
		db:                    db,
		zoneManager:           zoneManager,
		mapData:               mapData,
		playerPacketQueue:     shared.NewSafeQueue[[]byte](4096),
		mainServerPacketQueue: shared.NewSafeQueue[[]byte](4096),
		playerLoginQueue:      shared.NewSafeQueue[uint32](4096),
	}, nil
}

func (z *Zone) Start() error {
	z.logger.Info("Starting zone", shared.Field{Key: "mapId", Value: z.mapId})
	z.isRunning.Store(true)
	for z.isRunning.Load() {
		// TODO: game loop
	}

	z.logger.Info("Zone stopped", shared.Field{Key: "mapId", Value: z.mapId})
	return nil
}

func (z *Zone) EnqueuePlayerPacket(packet []byte) bool {
	return z.playerPacketQueue.Enqueue(packet)
}

func (z *Zone) EnqueueMainServerPacket(packet []byte) bool {
	return z.mainServerPacketQueue.Enqueue(packet)
}

func (z *Zone) EnqueuePlayerLogin(pcId uint32) bool {
	return z.playerLoginQueue.Enqueue(pcId)
}

func (z *Zone) Stop() {
	z.isRunning.Store(false)
}
