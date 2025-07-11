package webserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/db"
)

type Server struct {
	cfg    *config.EnvVars
	db     db.DBService
	logger shared.Logger
}

func NewServer(cfg *config.EnvVars, db db.DBService, logger shared.Logger) *http.Server {
	NewServer := &Server{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.cfg.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
