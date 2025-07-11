package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/db"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologFileLogger("web-server", "logs", cfg.GetZerologLevel())
	defer func(l shared.Logger) {
		_ = l.Close()
	}(logger)
	logger.Info("Starting Web Server service...")

	db, err := db.NewDbService(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
		os.Exit(1)
	}

	server := webserver.NewServer(cfg, db, logger)
	go func() {
		if err := server.Start(); err != nil {
			logger.Error("Server failed to start", shared.Field{Key: "error", Value: err})
			os.Exit(1)
		}
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info("Shutting down Web Server service...")
	_ = server.Stop(context.Background())
	_ = db.Close()
}
