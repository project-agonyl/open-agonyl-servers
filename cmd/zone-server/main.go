package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/zoneserver/db"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologFileLogger(fmt.Sprintf("zone-server-%d", cfg.ServerId), "logs", cfg.GetZerologLevel())
	defer func(l shared.Logger) {
		_ = l.Close()
	}(logger)
	logger.Info(fmt.Sprintf("Starting Zone Server %d service...", cfg.ServerId))
	db, err := db.NewDbService(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	cacheService := shared.NewRedisCacheService(cfg.CacheServerAddr, cfg.CacheServerPassword, cfg.CacheTlsEnabled)

	serialNumberGenerator := shared.NewSerialNumberGenerator(
		db.GetDB(),
		cacheService.(*redis.Client),
		fmt.Sprintf("zone-server-%d", cfg.ServerId),
	)

	players := zoneserver.NewPlayers()

	zoneManager := zoneserver.NewZoneManager(cfg, db, logger, cacheService.(*redis.Client), serialNumberGenerator, players)
	go func(z *zoneserver.ZoneManager) {
		err := z.Start()
		if err != nil {
			logger.Error("Failed to start zone manager", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(zoneManager)

	mainServerClient := zoneserver.NewMainServerClient(
		cfg.ServerId,
		cfg.MainServerIpAddress+":"+cfg.MainServerPort,
		logger,
		players,
		zoneManager,
		db,
	)
	go func(c *zoneserver.MainServerClient) {
		c.Start()
	}(mainServerClient)

	server := zoneserver.NewServer(cfg, db, logger, mainServerClient, players, zoneManager)
	go func(s *zoneserver.Server) {
		err := s.Start()
		if err != nil {
			logger.Error("Failed to start account server", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info(fmt.Sprintf("Shutting down Zone Server %d service...", cfg.ServerId))
	server.Stop()
	mainServerClient.Stop()
	_ = cacheService.Close()
	zoneManager.Stop()
	_ = db.Close()
}
