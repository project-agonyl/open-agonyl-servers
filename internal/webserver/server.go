package webserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/db"
)

type Server struct {
	cfg        *config.EnvVars
	db         db.DBService
	logger     shared.Logger
	httpServer *http.Server
}

func NewServer(cfg *config.EnvVars, db db.DBService, logger shared.Logger) *Server {
	server := &Server{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}

	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", server.cfg.Port),
		Handler:      server.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) Start() error {
	s.logger.Info("Starting web server", shared.Field{Key: "port", Value: s.cfg.Port})
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping web server")
	return s.httpServer.Shutdown(ctx)
}
