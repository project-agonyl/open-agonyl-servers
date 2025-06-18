package config

import (
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

type EnvVars struct {
	Port        string
	Environment string
	LogLevel    string
	DatabaseURL string
}

func New() *EnvVars {
	if _, ok := os.LookupEnv("PORT"); !ok {
		err := os.Setenv("PORT", "8080")
		if err != nil {
			slog.Info("Could not set default PORT!")
		}
	}

	if _, ok := os.LookupEnv("ENVIRONMENT"); !ok {
		err := os.Setenv("ENVIRONMENT", "production")
		if err != nil {
			slog.Info("Could not set default ENVIRONMENT!")
		}
	}

	if _, ok := os.LookupEnv("LOG_LEVEL"); !ok {
		err := os.Setenv("LOG_LEVEL", "info")
		if err != nil {
			slog.Info("Could not set default LOG_LEVEL!")
		}
	}

	if _, ok := os.LookupEnv("DATABASE_URL"); !ok {
		err := os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/agonyl?sslmode=disable")
		if err != nil {
			slog.Info("Could not set default DATABASE_URL!")
		}
	}

	return &EnvVars{
		Port:        os.Getenv("PORT"),
		Environment: os.Getenv("ENVIRONMENT"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func (e *EnvVars) GetZerologLevel() zerolog.Level {
	switch e.LogLevel {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
