# Project Structure and Utilities

## Directory Structure

### Root Level Directories

#### `/cmd/` - Application Entry Points
- **Purpose**: Contains main entry points for each server
- **Pattern**: Each server has its own subdirectory with `main.go`
- **Structure**:
  ```
  cmd/
  ├── account-server/main.go    # Account server entry point
  ├── gate-server/main.go       # Gate server entry point
  ├── login-server/main.go      # Login server entry point
  ├── main-server/main.go       # Main server entry point
  ├── migrate/main.go           # Database migration tool
  ├── web-server/main.go        # Web server entry point
  └── zone-server/main.go       # Zone server entry point (to be implemented)
  ```

#### `/internal/` - Private Application Code
- **Purpose**: Contains all internal application logic
- **Pattern**: Each server has its own package with consistent structure
- **Structure**:
  ```
  internal/
  ├── accountserver/           # Account server implementation
  ├── gateserver/             # Gate server implementation
  ├── loginserver/            # Login server implementation
  ├── mainserver/             # Main server implementation
  ├── shared/                 # Shared utilities and components
  ├── utils/                  # General utility functions
  ├── webserver/              # Web server implementation
  └── zoneserver/             # Zone server implementation (to be implemented)
  ```

#### `/build/` - Build and Deployment
- **Purpose**: Docker configurations and deployment scripts
- **Contents**:
  ```
  build/
  ├── docker-compose.yml      # Multi-service orchestration
  ├── env.example             # Environment configuration template
  ├── Dockerfile.*            # Docker files for each service
  └── README.md              # Build and deployment documentation
  ```

#### `/scripts/` - Build and Run Scripts
- **Purpose**: Cross-platform build and run scripts
- **Pattern**: Both `.sh` and `.bat` versions for each script
- **Contents**:
  ```
  scripts/
  ├── build.sh/bat           # Build individual server
  ├── build-all.sh/bat       # Build all servers
  ├── run.sh/bat             # Run individual server
  ├── run-all.sh/bat         # Run all servers
  ├── test.sh/bat            # Run tests
  └── migrate.sh/bat         # Run database migrations
  ```

#### `/logs/` - Application Logs
- **Purpose**: Centralized log storage
- **Pattern**: Each server writes to its own log file
- **Contents**: Generated at runtime, not in version control

### Server Package Structure

Each server follows a consistent internal structure:

```
internal/[servername]/
├── config/
│   └── config.go            # Environment-based configuration
├── db/
│   └── db.go                # Database operations with sqlx
├── server.go                # Main server implementation
├── server_session.go        # TCP session handling
├── player.go                # Player data structures
├── players.go               # Player collection management
└── [server]_client.go       # Client for other server communication
```

### Shared Components (`internal/shared/`)

#### `/shared/network/` - Network Infrastructure
- **Purpose**: Common networking components
- **Contents**:
  ```
  network/
  ├── tcp_server.go          # Base TCP server implementation
  └── tcp_session.go         # TCP session interface
  ```

#### `/shared/messages/` - Protocol Messages
- **Purpose**: Game protocol message definitions
- **Contents**:
  ```
  messages/
  ├── common.go              # Common message structures
  ├── account_server.go      # Account server messages
  ├── game_client.go         # Game client messages
  ├── gate_server.go         # Gate server messages
  ├── login_server.go        # Login server messages
  ├── main_server.go         # Main server messages
  └── protocol/
      └── protocol.go        # Protocol constants
  ```

#### `/shared/data/` - Game Data Structures
- **Purpose**: Binary file loading and game data structures
- **Contents**:
  ```
  data/
  ├── item.go                # Item data loading
  ├── map.go                 # Map data loading
  ├── npc.go                 # NPC data loading
  └── spawn.go               # Spawn data loading
  ```

#### `/shared/constants/` - System Constants
- **Purpose**: Shared constants and error codes
- **Contents**:
  ```
  constants/
  ├── constants.go           # General constants
  └── error.go              # Error code definitions
  ```

#### `/shared/crypto/` - Encryption
- **Purpose**: Custom encryption implementation
- **Contents**:
  ```
  crypto/
  └── crypto.go             # Packet encryption/decryption
  ```

#### `/shared/helpers/` - Helper Functions
- **Purpose**: Reusable helper functions
- **Contents**:
  ```
  helpers/
  ├── login.go              # Login-related helpers
  └── settings.go           # Settings management helpers
  ```

### Utility Components (`internal/utils/`)

#### `/utils/` - General Utilities
- **Purpose**: Cross-cutting utility functions
- **Contents**:
  ```
  utils/
  ├── bytes.go              # Byte manipulation utilities
  ├── performance.go        # Performance monitoring utilities
  ├── string.go             # String manipulation utilities
  └── ull.go               # Unsigned long long utilities
  ```

## Available Utilities and Helpers

### Thread-Safe Data Structures

#### `SafeMap[K, V]` - Thread-Safe Map
```go
// Usage example
players := shared.NewSafeMap[uint32, *Player]()
players.Store(pcId, player)
player, exists := players.Get(pcId)
players.Delete(pcId)
players.Range(func(k uint32, v *Player) bool {
    // Process each player
    return true
})
```

**Methods**:
- `Store(k K, v V)` - Store a key-value pair
- `Get(k K) (V, bool)` - Retrieve a value
- `Delete(k K)` - Remove a key-value pair
- `Range(f func(k K, v V) bool)` - Iterate over all entries
- `Has(k K) bool` - Check if key exists
- `Len() int` - Get map size

#### `SafeSet[T]` - Thread-Safe Set
```go
// Usage example
onlinePlayers := shared.NewSafeSet[uint32]()
onlinePlayers.Add(pcId)
if onlinePlayers.Contains(pcId) {
    // Player is online
}
onlinePlayers.Remove(pcId)
```

**Methods**:
- `Add(value T)` - Add element to set
- `Remove(value T)` - Remove element from set
- `Contains(value T) bool` - Check if element exists
- `Size() int` - Get set size
- `Range(f func(value T) bool)` - Iterate over elements
- `List() []T` - Get all elements as slice

### Logging System

#### `Logger` Interface
```go
// Usage example
logger := shared.NewZerologFileLogger("zone-server", "logs", zerolog.InfoLevel)
logger.Info("Player connected", 
    shared.Field{Key: "pcId", Value: pcId},
    shared.Field{Key: "username", Value: username},
)
logger.Error("Database error", shared.Field{Key: "error", Value: err})
```

**Features**:
- Structured logging with fields
- File rotation with daily logs
- Multiple log levels (Debug, Info, Warn, Error)
- Context-aware logging with `With()` method

### Database Utilities

#### `sqlx` + `squirrel` Pattern
```go
// Usage example
func (d *DBService) GetPlayerData(pcId uint32) (*PlayerData, error) {
    query := squirrel.Select("*").
        From("players").
        Where(squirrel.Eq{"pc_id": pcId}).
        Limit(1)
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    var player PlayerData
    err = d.db.Get(&player, sql, args...)
    return &player, err
}
```

**Features**:
- Type-safe query building with squirrel
- Connection pooling with sqlx
- Manual struct mapping for performance
- Parameterized queries for security

### Network Utilities

#### `TCPServer` Base Implementation
```go
// Usage example
server := &Server{
    TCPServer: network.TCPServer{
        Addr:         cfg.IpAddress + ":" + cfg.Port,
        Name:         "zone-server",
        UidGenerator: shared.NewUidGenerator(0),
        Logger:       logger,
        Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
    },
}
```

**Features**:
- Automatic session management
- Graceful shutdown handling
- Connection pooling
- Logging integration

### Binary File Loading

#### Map Data Loading
```go
// Usage example
mapData, err := data.LoadMapData("maps/map_1.dat")
if err != nil {
    return err
}
// mapData contains: Id, Name, WarpData, NavigationMesh
```

#### NPC Data Loading
```go
// Usage example
npcData, err := data.LoadNPCData("npcs/npc_100.n_ndt")
if err != nil {
    return err
}
// npcData contains: Name, Id, HP, Attacks, etc.
```

#### Spawn Data Loading
```go
// Usage example
spawns, err := data.LoadNPCSpawnData("spawns/spawn_1.n_ndt")
if err != nil {
    return err
}
// spawns contains: Id, X, Y, Orientation, SpawnStep
```

### String Utilities

#### `ReadStringFromBytes`
```go
// Usage example
name := utils.ReadStringFromBytes(buffer)
// Handles null-terminated strings from binary data
```

#### `GenerateRandomString`
```go
// Usage example
token := utils.GenerateRandomString(32)
// Generates random alphanumeric string
```

### Performance Monitoring

#### `PerformanceMonitor`
```go
// Usage example
monitor := utils.NewPerformanceMonitor()
monitor.Start()
// ... perform operation ...
monitor.Stop()
elapsed := monitor.ElapsedMilliseconds()
```

### Cache Service

#### Redis/Valkey Integration
```go
// Usage example
cache := shared.NewRedisCacheService(addr, password, tlsEnabled)
err := cache.Set(ctx, "key", "value", time.Hour).Err()
value, err := cache.Get(ctx, "key").Result()
```

### Crypto Utilities

#### Custom Encryption
```go
// Usage example
crypto := crypto.NewCrypto562(dynamicKey)
encrypted := crypto.Encrypt(packet)
decrypted := crypto.Decrypt(packet)
```

### ID Generation

#### `UidGenerator`
```go
// Usage example
generator := shared.NewUidGenerator(serverId)
sessionId := generator.Generate()
```

#### `SerialNumberGenerator`
```go
// Usage example
serialGen := shared.NewSerialNumberGenerator(db, cache, prefix)
serial := serialGen.Generate()
```

## Configuration Pattern

### Environment-Based Configuration
```go
// Usage example
type EnvVars struct {
    Port        string
    IpAddress   string
    ServerId    byte
    DatabaseURL string
    // ... other fields
}

func New() *EnvVars {
    // Set defaults if not present
    if _, ok := os.LookupEnv("PORT"); !ok {
        os.Setenv("PORT", "default_port")
    }
    
    // Parse environment variables
    port := os.Getenv("PORT")
    // ... parse other fields
    
    return &EnvVars{
        Port: port,
        // ... other fields
    }
}
```

## Reusable Patterns

### Server Initialization Pattern
```go
func main() {
    cfg := config.New()
    logger := shared.NewZerologFileLogger("server-name", "logs", cfg.GetZerologLevel())
    defer logger.Close()
    
    db, err := db.NewDbService(cfg.DatabaseURL, logger)
    if err != nil {
        logger.Error("Failed to create db service", shared.Field{Key: "error", Value: err})
        os.Exit(1)
    }
    
    server := NewServer(cfg, db, logger)
    go server.Start()
    
    // Graceful shutdown
    interruptChan := make(chan os.Signal, 1)
    signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
    <-interruptChan
    
    server.Stop()
    db.Close()
}
```

### Session Management Pattern
```go
type serverSession struct {
    server   *Server
    conn     net.Conn
    id       uint32
    sendChan chan []byte
    done     chan struct{}
    wg       sync.WaitGroup
}

func (s *serverSession) Handle() {
    defer func() {
        close(s.done)
        s.wg.Wait()
        s.server.RemoveSession(s.id)
    }()
    
    // Handle incoming packets
    for {
        // Read packet
        // Process packet
        // Send response
    }
}
```

### Database Query Pattern
```go
func (d *DBService) GetData(id uint32) (*Data, error) {
    query := squirrel.Select("*").
        From("table").
        Where(squirrel.Eq{"id": id})
    
    sql, args, err := query.ToSql()
    if err != nil {
        return nil, err
    }
    
    var data Data
    err = d.db.Get(&data, sql, args...)
    return &data, err
}
```

This comprehensive structure provides a solid foundation for building the zone server while maintaining consistency with existing patterns and avoiding code duplication. 