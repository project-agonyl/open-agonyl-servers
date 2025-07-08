package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/accountserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologFileLogger("account-server", "logs", cfg.GetZerologLevel())
	defer func(l shared.Logger) {
		_ = l.Close()
	}(logger)
	logger.Info("Starting Account Server service...")
	db, err := db.NewDbService(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	players := accountserver.NewPlayers()
	mainServerClient := accountserver.NewMainServerClient(
		cfg.ServerId,
		cfg.MainServerIpAddress+":"+cfg.MainServerPort,
		logger,
		players,
	)
	go func(c *accountserver.MainServerClient) {
		c.Start()
	}(mainServerClient)

	cacheService := shared.NewRedisCacheService(cfg.CacheServerAddr, cfg.CacheServerPassword, cfg.CacheTlsEnabled)
	serialNumberGenerator := shared.NewSerialNumberGenerator(
		db.GetDB(),
		cacheService.(*redis.Client),
		fmt.Sprintf("account-server-%d", cfg.ServerId),
	)
	server := accountserver.NewServer(cfg, db, logger, players, mainServerClient, serialNumberGenerator)
	go func(s *accountserver.Server) {
		err := s.Start()
		if err != nil {
			logger.Error("Failed to start account server", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info("Shutting down Account Server service...")
	server.Stop()
	mainServerClient.Stop()
	_ = cacheService.Close()
	_ = db.Close()
}
