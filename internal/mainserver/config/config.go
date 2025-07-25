package config

import (
	"log/slog"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

type EnvVars struct {
	Port                string
	IpAddress           string
	Environment         string
	LogLevel            string
	DatabaseURL         string
	CacheServerAddr     string
	CacheServerPassword string
	CacheTlsEnabled     bool
	CacheKeyPrefix      string
}

func New() *EnvVars {
	if _, ok := os.LookupEnv("PORT"); !ok {
		err := os.Setenv("PORT", "5555")
		if err != nil {
			slog.Info("Could not set default PORT!")
		}
	}

	if _, ok := os.LookupEnv("IP_ADDRESS"); !ok {
		err := os.Setenv("IP_ADDRESS", "127.0.0.1")
		if err != nil {
			slog.Info("Could not set default IP_ADDRESS!")
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

	if _, ok := os.LookupEnv("CACHE_SERVER_ADDR"); !ok {
		err := os.Setenv("CACHE_SERVER_ADDR", "localhost:6379")
		if err != nil {
			slog.Info("Could not set default CACHE_SERVER_ADDR!")
		}
	}

	if _, ok := os.LookupEnv("CACHE_SERVER_PASSWORD"); !ok {
		err := os.Setenv("CACHE_SERVER_PASSWORD", "")
		if err != nil {
			slog.Info("Could not set default CACHE_SERVER_PASSWORD!")
		}
	}

	if _, ok := os.LookupEnv("CACHE_TLS_ENABLED"); !ok {
		err := os.Setenv("CACHE_TLS_ENABLED", "false")
		if err != nil {
			slog.Info("Could not set default CACHE_TLS_ENABLED!")
		}
	}

	cacheTlsEnabled, err := strconv.ParseBool(os.Getenv("CACHE_TLS_ENABLED"))
	if err != nil {
		cacheTlsEnabled = false
	}

	if _, ok := os.LookupEnv("CACHE_KEY_PREFIX"); !ok {
		err := os.Setenv("CACHE_KEY_PREFIX", "agonyl:main-server:")
		if err != nil {
			slog.Info("Could not set default CACHE_KEY_PREFIX!")
		}
	}

	return &EnvVars{
		Port:                os.Getenv("PORT"),
		IpAddress:           os.Getenv("IP_ADDRESS"),
		Environment:         os.Getenv("ENVIRONMENT"),
		LogLevel:            os.Getenv("LOG_LEVEL"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		CacheServerAddr:     os.Getenv("CACHE_SERVER_ADDR"),
		CacheServerPassword: os.Getenv("CACHE_SERVER_PASSWORD"),
		CacheTlsEnabled:     cacheTlsEnabled,
		CacheKeyPrefix:      os.Getenv("CACHE_KEY_PREFIX"),
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
