package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
)

type ZoneServerInfo struct {
	ID   int
	IP   string
	Port int
}

type EnvVars struct {
	Port                 string
	IpAddress            string
	ServerName           string
	ServerId             byte
	Environment          string
	LogLevel             string
	DatabaseURL          string
	CacheServerAddr      string
	CacheServerPassword  string
	CacheTlsEnabled      bool
	CacheKeyPrefix       string
	IsTestMode           bool
	LoginServerIpAddress string
	LoginServerPort      string
	DynamicKey           int
	ZoneServers          []ZoneServerInfo
}

func New() *EnvVars {
	if _, ok := os.LookupEnv("PORT"); !ok {
		err := os.Setenv("PORT", "9860")
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

	if _, ok := os.LookupEnv("SERVER_ID"); !ok {
		err := os.Setenv("SERVER_ID", "0")
		if err != nil {
			slog.Info("Could not set default SERVER_ID!")
		}
	}

	serverId, err := strconv.ParseUint(os.Getenv("SERVER_ID"), 10, 8)
	if err != nil {
		serverId = 0
	}

	if _, ok := os.LookupEnv("ENVIRONMENT"); !ok {
		err := os.Setenv("ENVIRONMENT", "production")
		if err != nil {
			slog.Info("Could not set default ENVIRONMENT!")
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

	if _, ok := os.LookupEnv("LOGIN_SERVER_IP_ADDRESS"); !ok {
		err := os.Setenv("LOGIN_SERVER_IP_ADDRESS", "127.0.0.1")
		if err != nil {
			slog.Info("Could not set default LOGIN_SERVER_IP_ADDRESS!")
		}
	}

	if _, ok := os.LookupEnv("LOGIN_SERVER_PORT"); !ok {
		err := os.Setenv("LOGIN_SERVER_PORT", "3210")
		if err != nil {
			slog.Info("Could not set default LOGIN_SERVER_PORT!")
		}
	}

	if _, ok := os.LookupEnv("DYNAMIC_KEY"); !ok {
		err := os.Setenv("DYNAMIC_KEY", "0x04C478BD")
		if err != nil {
			slog.Info("Could not set default DYNAMIC_KEY!")
		}
	}

	var dynamicKey int
	d, err := strconv.ParseInt(os.Getenv("DYNAMIC_KEY"), 0, 32)
	if err != nil {
		dynamicKey = 0x04C478BD
	} else {
		dynamicKey = int(d)
	}

	if _, ok := os.LookupEnv("SERVER_NAME"); !ok {
		err := os.Setenv("SERVER_NAME", "Agonyl")
		if err != nil {
			slog.Info("Could not set default SERVER_NAME!")
		}
	}

	if _, ok := os.LookupEnv("ZONE_SERVER_COUNT"); !ok {
		err := os.Setenv("ZONE_SERVER_COUNT", "3")
		if err != nil {
			slog.Info("Could not set default ZONE_SERVER_COUNT!")
		}
	}

	defaultZoneServers := []ZoneServerInfo{
		{
			ID:   255,
			IP:   "127.0.0.1",
			Port: 5589,
		},
		{
			ID:   0,
			IP:   "127.0.0.1",
			Port: 7568,
		},
		{
			ID:   3,
			IP:   "127.0.0.1",
			Port: 6699,
		},
	}

	zoneServerCount, err := strconv.Atoi(os.Getenv("ZONE_SERVER_COUNT"))
	if err != nil {
		zoneServerCount = len(defaultZoneServers)
	}

	zoneServers := make([]ZoneServerInfo, 0)
	for i := range make([]int, zoneServerCount) {
		zoneServerIP := os.Getenv(fmt.Sprintf("ZONE_SERVER_IP_%d", i+1))
		if zoneServerIP == "" {
			continue
		}

		zoneServerPort, err := strconv.Atoi(os.Getenv(fmt.Sprintf("ZONE_SERVER_PORT_%d", i+1)))
		if err != nil {
			continue
		}

		zoneServerID, err := strconv.Atoi(os.Getenv(fmt.Sprintf("ZONE_SERVER_ID_%d", i+1)))
		if err != nil {
			continue
		}

		zoneServers = append(zoneServers, ZoneServerInfo{
			ID:   zoneServerID,
			IP:   zoneServerIP,
			Port: zoneServerPort,
		})
	}

	if len(zoneServers) == 0 {
		zoneServers = defaultZoneServers
	}

	return &EnvVars{
		Port:                 os.Getenv("PORT"),
		Environment:          os.Getenv("ENVIRONMENT"),
		IpAddress:            os.Getenv("IP_ADDRESS"),
		ServerId:             byte(serverId),
		LogLevel:             os.Getenv("LOG_LEVEL"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		CacheServerAddr:      os.Getenv("CACHE_SERVER_ADDR"),
		CacheServerPassword:  os.Getenv("CACHE_SERVER_PASSWORD"),
		CacheTlsEnabled:      cacheTlsEnabled,
		CacheKeyPrefix:       os.Getenv("CACHE_KEY_PREFIX"),
		IsTestMode:           isTestMode,
		LoginServerIpAddress: os.Getenv("LOGIN_SERVER_IP_ADDRESS"),
		LoginServerPort:      os.Getenv("LOGIN_SERVER_PORT"),
		DynamicKey:           dynamicKey,
		ServerName:           os.Getenv("SERVER_NAME"),
		ZoneServers:          zoneServers,
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

func (e *EnvVars) GetLoginServerPort() int {
	port, err := strconv.Atoi(e.LoginServerPort)
	if err != nil {
		return 3210
	}

	return port
}

func (e *EnvVars) GetServerPort() int {
	port, err := strconv.Atoi(e.Port)
	if err != nil {
		return 9860
	}

	return port
}
