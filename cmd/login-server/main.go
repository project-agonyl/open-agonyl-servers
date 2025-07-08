package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/loginserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologFileLogger("login-server", "logs", cfg.GetZerologLevel())
	defer func(l shared.Logger) {
		_ = l.Close()
	}(logger)
	logger.Info(
		"Starting Login Server service...",
		shared.Field{Key: "autoCreateAccount", Value: cfg.AutoCreateAccount},
		shared.Field{Key: "isTestMode", Value: cfg.IsTestMode},
	)
	db, err := db.NewDbService(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	cache := shared.NewRedisCacheService(cfg.CacheServerAddr, cfg.CacheServerPassword, cfg.CacheTlsEnabled)
	_, err = cache.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("Failed to ping cache service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	broker := loginserver.NewBroker(fmt.Sprintf(":%s", cfg.BrokerPort), logger, cache)
	go func(ls *loginserver.Broker) {
		err := ls.Start()
		if err != nil {
			logger.Error("Failed to start broker", shared.Field{Key: "error", Value: err})
			panic(err)
		}

	}(broker)

	server := loginserver.NewServer(fmt.Sprintf(":%s", cfg.Port), logger, cache, db, broker, cfg)
	go func(ls *loginserver.Server) {
		err := ls.Start()
		if err != nil {
			logger.Error("Failed to start login server", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info("Shutting down Login Server service...")
	broker.Stop()
	server.Stop()
	_ = cache.Close()
	_ = db.Close()
}
