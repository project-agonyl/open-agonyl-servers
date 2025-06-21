package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/gateserver/db"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/crypto"
	"github.com/rs/zerolog"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologLogger(zerolog.New(os.Stdout), "gate-server", cfg.GetZerologLevel())
	logger.Info("Starting Gate Server service...")
	db, err := db.NewDbService(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	lsClient := gateserver.NewLoginServerClient(cfg, logger)
	go lsClient.Start()
	crypt := crypto.NewCrypto562(cfg.DynamicKey)
	players := gateserver.NewPlayers()
	zsClients := gateserver.NewZoneServerClients(cfg, crypt, players, logger)
	go zsClients.Start()
	server := gateserver.NewServer(logger, cfg, players, zsClients, lsClient, crypt, db)
	go func(s *gateserver.Server) {
		err := s.Start()
		if err != nil {
			logger.Error("Failed to start gate server", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info("Shutting down Gate Server service...")
	server.Stop()
	zsClients.Stop()
	lsClient.Stop()
	_ = db.Close()
}
