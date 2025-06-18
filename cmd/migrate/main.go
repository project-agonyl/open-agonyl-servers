package main

import (
	"database/sql"
	"embed"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/project-agonyl/open-agonyl-servers/cmd/migrate/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/rs/zerolog"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	cfg := config.New()
	logger := shared.NewZerologLogger(
		zerolog.New(os.Stdout), "migrate", cfg.GetZerologLevel())

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("database connection error", shared.Field{
			Key:   "error",
			Value: err,
		})
		os.Exit(1)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error("postgres driver error", shared.Field{
			Key:   "error",
			Value: err,
		})
		os.Exit(1)
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Error("migration source error", shared.Field{
			Key:   "error",
			Value: err,
		})
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		logger.Error("migrate init error", shared.Field{
			Key:   "error",
			Value: err,
		})
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		logger.Error("migrate up error", shared.Field{
			Key:   "error",
			Value: err,
		})
		os.Exit(1)
	}

	logger.Info("migration completed successfully")
}
