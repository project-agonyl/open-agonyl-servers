package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/project-agonyl/open-agonyl-servers/internal/mainserver"
	"github.com/project-agonyl/open-agonyl-servers/internal/mainserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

func main() {
	cfg := config.New()
	logger := shared.NewZerologFileLogger("main-server", "logs", cfg.GetZerologLevel())
	defer func(l shared.Logger) {
		_ = l.Close()
	}(logger)
	logger.Info("Starting Main Server service...")
	server := mainserver.NewServer(cfg, logger)
	go func(s *mainserver.Server) {
		err := s.Start()
		if err != nil {
			logger.Error("Failed to start main server", shared.Field{Key: "error", Value: err})
			panic(err)
		}
	}(server)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	<-interruptChan
	logger.Info("Shutting down Main Server service...")
	server.Stop()
}
