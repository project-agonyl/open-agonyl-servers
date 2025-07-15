# Technical Context

## Technology Stack

### Core Technologies
- **Language**: Go 1.21+
- **Database**: PostgreSQL with sqlx for database operations
- **Query Builder**: Squirrel for type-safe SQL query building
- **Cache**: Redis/Valkey for session and temporary data
- **Networking**: Custom TCP protocol implementation
- **Logging**: Zerolog for structured logging
- **Configuration**: Environment variables with defaults

### Database Layer
- **No ORM**: Direct SQL queries using sqlx
- **Query Building**: Squirrel for dynamic query construction
- **Connection Pooling**: sqlx built-in connection pooling
- **Migration System**: Custom migration scripts in SQL

### Data Storage
- **Game Data**: Binary files for maps, NPCs, items, spawns
- **Player Data**: PostgreSQL for persistent player information
- **Session Data**: Redis for temporary session storage
- **Cache Data**: Redis for frequently accessed data

### Network Protocol
- **Custom Protocol**: Binary protocol with specific message formats
- **TCP Server**: Custom implementation with session management
- **Encryption**: Custom crypto implementation for packet security
- **Message Types**: Structured message headers with protocol IDs

## Development Environment

### Build System
- **Go Modules**: Modern Go dependency management
- **Docker**: Containerized development and deployment
- **Scripts**: Build and run scripts for each server component

### Testing
- **Unit Tests**: Go testing framework
- **Integration Tests**: Database and network testing
- **Performance Tests**: Load testing for server components

### Deployment
- **Docker Compose**: Multi-service orchestration
- **Environment Configuration**: Environment-specific settings
- **Logging**: File-based logging with rotation

## Performance Characteristics

### Database Performance
- **Connection Pooling**: Efficient database connection management
- **Query Optimization**: Manual query optimization without ORM overhead
- **Indexing**: Strategic database indexing for game queries
- **Caching**: Redis caching for frequently accessed data

### Network Performance
- **TCP Optimization**: TCP_NODELAY for low latency
- **Packet Processing**: Efficient binary packet handling
- **Concurrent Connections**: Goroutine-based connection handling
- **Memory Management**: Minimal allocations in network code

### Game Performance
- **Spatial Partitioning**: Efficient player/NPC location tracking
- **Event System**: Asynchronous event processing
- **Tick Rate**: Configurable simulation tick rates
- **Memory Pooling**: Reuse of frequently allocated objects

## Security Considerations

### Network Security
- **Custom Encryption**: Packet-level encryption
- **Session Management**: Secure session handling
- **Input Validation**: Packet validation and sanitization
- **Rate Limiting**: Protection against abuse

### Data Security
- **Password Hashing**: Secure password storage
- **SQL Injection Prevention**: Parameterized queries with sqlx
- **Access Control**: Role-based access control
- **Audit Logging**: Comprehensive security logging

## Scalability Features

### Horizontal Scaling
- **Zone-Based Architecture**: Independent zone servers
- **Load Balancing**: Distribution across multiple servers
- **State Management**: Efficient state synchronization
- **Database Sharding**: Potential for database partitioning

### Vertical Scaling
- **Resource Optimization**: Efficient memory and CPU usage
- **Connection Pooling**: Optimized database connections
- **Caching Strategy**: Multi-level caching approach
- **Background Processing**: Asynchronous task processing

## Monitoring and Observability

### Logging
- **Structured Logging**: JSON-formatted logs with fields
- **Log Levels**: Configurable logging verbosity
- **Log Rotation**: Automatic log file management
- **Context Tracking**: Request/response correlation

### Metrics
- **Performance Metrics**: Response times and throughput
- **Resource Metrics**: CPU, memory, and network usage
- **Business Metrics**: Player counts and game statistics
- **Error Tracking**: Error rates and types

## Development Workflow

### Code Organization
- **Microservices**: Separate packages for each server
- **Shared Code**: Common utilities in internal/shared
- **Configuration**: Environment-based configuration
- **Testing**: Comprehensive test coverage

### Version Control
- **Git Workflow**: Feature branch development
- **Code Review**: Peer review process
- **Continuous Integration**: Automated testing and building
- **Deployment**: Automated deployment pipeline 