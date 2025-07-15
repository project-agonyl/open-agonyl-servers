# System Patterns

## Architecture Patterns

### Microservices Architecture
The system follows a microservices pattern with specialized servers:
- Each server has a specific responsibility
- Servers communicate via TCP with custom protocols
- Shared data structures in `internal/shared/`
- Common networking patterns in `internal/shared/network/`

### TCP Server Pattern
All servers follow a consistent TCP server pattern:
```go
type Server struct {
    network.TCPServer
    // Server-specific fields
}

func NewServer(...) *Server {
    server := &Server{
        TCPServer: network.TCPServer{
            Addr:         addr,
            Name:         "server-name",
            UidGenerator: shared.NewUidGenerator(0),
            Logger:       logger,
            Sessions:     shared.NewSafeMap[uint32, network.TCPServerSession](),
        },
        // Initialize server-specific fields
    }
    server.NewSession = func(id uint32, conn net.Conn) network.TCPServerSession {
        session := newServerSession(id, conn)
        // Set server reference
        return session
    }
    return server
}
```

### Session Management Pattern
Each server implements session handling:
- Sessions are created per connection
- Each session has a unique ID
- Sessions handle packet processing
- Sessions manage player state

### Client-Server Communication Pattern
Servers communicate with each other using client patterns:
- `MainServerClient` for account server
- `ZoneServerClient` for gate server
- `LoginServerClient` for gate server
- Reconnection logic with exponential backoff

## Data Patterns

### Database Access Pattern
The system uses **sqlx** for database operations and **squirrel** for query building:
- No ORM is used - direct SQL queries with sqlx
- Squirrel provides type-safe query building
- Manual mapping between database rows and structs
- Connection pooling with sqlx

```go
// Example database pattern
type DBService struct {
    db     *sqlx.DB
    logger shared.Logger
}

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
    if err != nil {
        return nil, err
    }
    
    return &player, nil
}
```

### Binary File Loading
Game data is loaded from binary files:
- Map data: `LoadMapData(mapFilePath string)`
- NPC data: `LoadNPCData(npcFilePath string)`
- Spawn data: `LoadNPCSpawnData(spawnFilePath string)`
- Item data: `LoadIT0Items()`, `LoadIT1Items()`, etc.

### Safe Concurrent Data Structures
- `SafeMap` for thread-safe maps
- `SafeSet` for thread-safe sets
- Atomic operations for counters and flags

### Message Protocol Pattern
All messages follow a consistent structure:
```go
type MsgHead struct {
    Size     uint16
    Protocol uint16
    MsgHeadNoProtocol
}

type MsgHeadNoProtocol struct {
    Ctrl  byte
    Cmd   byte
    PcId  uint32
}
```

## Configuration Patterns

### Environment-Based Configuration
All servers use environment variables for configuration:
- Database URLs
- Server addresses and ports
- Cache settings
- Feature flags

### Default Value Pattern
Configuration provides sensible defaults:
```go
if _, ok := os.LookupEnv("PORT"); !ok {
    err := os.Setenv("PORT", "default_port")
}
```

## Logging Patterns

### Structured Logging
All servers use structured logging with fields:
```go
logger.Info("message", 
    shared.Field{Key: "key", Value: value},
    shared.Field{Key: "another", Value: anotherValue},
)
```

### File-Based Logging
Logs are written to files in the `logs/` directory with rotation.

## Error Handling Patterns

### Graceful Degradation
- Connection failures trigger reconnection
- Database errors are logged but don't crash
- Invalid packets are logged and ignored

### Resource Cleanup
- Proper connection closing
- Session cleanup on disconnect
- Database connection pooling with sqlx

## Security Patterns

### Encryption
- Custom crypto implementation for packet encryption
- Dynamic key generation
- Encrypted communication between servers

### Authentication
- Session-based authentication
- Token-based login system
- Account validation

## Performance Patterns

### Goroutine Management
- Each connection runs in its own goroutine
- Packet processing in separate goroutines
- Background tasks for maintenance

### Memory Management
- Buffered channels for communication
- Efficient binary data structures
- Minimal allocations in hot paths

### Database Optimization
- Connection pooling with sqlx
- Prepared statements for frequently used queries
- Efficient query building with squirrel
- Manual query optimization without ORM overhead 