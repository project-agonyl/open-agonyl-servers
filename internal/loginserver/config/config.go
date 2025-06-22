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
	Environment         string
	LogLevel            string
	DatabaseURL         string
	CacheServerAddr     string
	CacheServerPassword string
	CacheTlsEnabled     bool
	CacheKeyPrefix      string
	IsTestMode          bool
	BrokerPort          string
	AutoCreateAccount   bool
}

func New() *EnvVars {
	if _, ok := os.LookupEnv("PORT"); !ok {
		err := os.Setenv("PORT", "3550")
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
		err := os.Setenv("CACHE_KEY_PREFIX", "agonyl:login-server:")
		if err != nil {
			slog.Info("Could not set default CACHE_KEY_PREFIX!")
		}
	}

	if _, ok := os.LookupEnv("IS_TEST_MODE"); !ok {
		err := os.Setenv("IS_TEST_MODE", "false")
		if err != nil {
			slog.Info("Could not set default IS_TEST_MODE!")
		}
	}

	isTestMode, err := strconv.ParseBool(os.Getenv("IS_TEST_MODE"))
	if err != nil {
		isTestMode = false
	}

	if _, ok := os.LookupEnv("BROKER_PORT"); !ok {
		err := os.Setenv("BROKER_PORT", "3210")
		if err != nil {
			slog.Info("Could not set default BROKER_PORT!")
		}
	}

	if _, ok := os.LookupEnv("AUTO_CREATE_ACCOUNT"); !ok {
		err := os.Setenv("AUTO_CREATE_ACCOUNT", "false")
		if err != nil {
			slog.Info("Could not set default AUTO_CREATE_ACCOUNT!")
		}
	}

	autoCreateAccount, err := strconv.ParseBool(os.Getenv("AUTO_CREATE_ACCOUNT"))
	if err != nil {
		autoCreateAccount = false
	}

	return &EnvVars{
		Port:                os.Getenv("PORT"),
		Environment:         os.Getenv("ENVIRONMENT"),
		LogLevel:            os.Getenv("LOG_LEVEL"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		CacheServerAddr:     os.Getenv("CACHE_SERVER_ADDR"),
		CacheServerPassword: os.Getenv("CACHE_SERVER_PASSWORD"),
		CacheTlsEnabled:     cacheTlsEnabled,
		CacheKeyPrefix:      os.Getenv("CACHE_KEY_PREFIX"),
		IsTestMode:          isTestMode,
		BrokerPort:          os.Getenv("BROKER_PORT"),
		AutoCreateAccount:   autoCreateAccount,
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
