package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

type EnvVars struct {
	Port                          string
	Environment                   string
	LogLevel                      string
	DatabaseURL                   string
	CacheServerAddr               string
	CacheServerPassword           string
	CacheTlsEnabled               bool
	CacheKeyPrefix                string
	ServerName                    string
	JwtSecret                     string
	JwtExpiry                     int
	SessionCookieName             string
	IsAccountVerificationRequired bool
}

func New() *EnvVars {
	if _, ok := os.LookupEnv("PORT"); !ok {
		err := os.Setenv("PORT", "80")
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
		err := os.Setenv("CACHE_KEY_PREFIX", "agonyl:account-server:")
		if err != nil {
			slog.Info("Could not set default CACHE_KEY_PREFIX!")
		}
	}

	if _, ok := os.LookupEnv("SERVER_NAME"); !ok {
		err := os.Setenv("SERVER_NAME", "A3 Agonyl")
		if err != nil {
			slog.Info("Could not set default SERVER_NAME!")
		}
	}

	if _, ok := os.LookupEnv("JWT_SECRET"); !ok {
		err := os.Setenv("JWT_SECRET", "d49063076fa02e28ba37752964780754de6ebec647e52493b036a5495d245c21")
		if err != nil {
			slog.Info("Could not set default JWT_SECRET!")
		}
	}

	if _, ok := os.LookupEnv("JWT_EXPIRY"); !ok {
		err := os.Setenv("JWT_EXPIRY", "1h")
		if err != nil {
			slog.Info("Could not set default JWT_EXPIRY!")
		}
	}

	jwtExpiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if err != nil {
		jwtExpiry = 7 * 24 * 60 * 60 // 7 days
	}

	if _, ok := os.LookupEnv("SESSION_COOKIE_NAME"); !ok {
		err := os.Setenv("SESSION_COOKIE_NAME", "agonyl_session_token")
		if err != nil {
			slog.Info("Could not set default SESSION_COOKIE_NAME!")
		}
	}

	if _, ok := os.LookupEnv("IS_ACCOUNT_VERIFICATION_REQUIRED"); !ok {
		err := os.Setenv("IS_ACCOUNT_VERIFICATION_REQUIRED", "false")
		if err != nil {
			slog.Info("Could not set default IS_ACCOUNT_VERIFICATION_REQUIRED!")
		}
	}

	isAccountVerificationRequired, err := strconv.ParseBool(os.Getenv("IS_ACCOUNT_VERIFICATION_REQUIRED"))
	if err != nil {
		isAccountVerificationRequired = false
	}

	return &EnvVars{
		Port:                          os.Getenv("PORT"),
		Environment:                   os.Getenv("ENVIRONMENT"),
		LogLevel:                      os.Getenv("LOG_LEVEL"),
		DatabaseURL:                   os.Getenv("DATABASE_URL"),
		CacheServerAddr:               os.Getenv("CACHE_SERVER_ADDR"),
		CacheServerPassword:           os.Getenv("CACHE_SERVER_PASSWORD"),
		CacheTlsEnabled:               cacheTlsEnabled,
		CacheKeyPrefix:                os.Getenv("CACHE_KEY_PREFIX"),
		ServerName:                    os.Getenv("SERVER_NAME"),
		JwtSecret:                     os.Getenv("JWT_SECRET"),
		JwtExpiry:                     jwtExpiry,
		SessionCookieName:             os.Getenv("SESSION_COOKIE_NAME"),
		IsAccountVerificationRequired: isAccountVerificationRequired,
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
