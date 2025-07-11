# Agonyl MMO Game Server - Docker Compose Setup

This directory contains the Docker Compose configuration for running the complete Agonyl MMO game server stack.

## Services

The stack includes the following services in startup order:

1. **Valkey** (Redis-compatible) - Cache and session storage
2. **PostgreSQL** - Primary database
3. **Adminer** - Database management interface
4. **Migrate** - Database migration service
5. **Main Server** - Core game server
6. **Account Server** - Character and account management
7. **Login Server** - Authentication service
8. **Gate Server** - Game client gateway

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- At least 4GB RAM available
- Ports 3210, 3550, 5432, 5555, 5589, 6379, 8080, 9860 available

## Quick Start

1. **Navigate to the build directory:**
   ```bash
   cd build
   ```

2. **Start all services:**
   ```bash
   docker-compose up -d
   ```

3. **Check service status:**
   ```bash
   docker-compose ps
   ```

4. **View logs:**
   ```bash
   # All services
   docker-compose logs -f
   
   # Specific service
   docker-compose logs -f main-server
   ```

## Service URLs

- **Adminer (Database Management):** http://localhost:8080
  - Server: `postgres`
  - Username: `postgres`
  - Password: `postgres`
  - Database: `agonyl`

- **Game Server Endpoints:**
  - Main Server: `localhost:5555`
  - Account Server: `localhost:5589`
  - Login Server: `localhost:3550`
  - Gate Server: `localhost:9860`
  - Login Broker: `localhost:3210`

## Configuration

### Environment Variables

Each service uses environment variables for configuration. Key variables include:

- `DATABASE_URL`: PostgreSQL connection string
- `CACHE_SERVER_ADDR`: Valkey/Redis server address
- `LOG_LEVEL`: Logging level (trace, debug, info, warn, error)
- `ENVIRONMENT`: Runtime environment (production, development)

### Volumes

- `postgres_data`: PostgreSQL data persistence
- `valkey_data`: Valkey data persistence
- `./logs`: Application logs (mounted to host)

## Development

### Building Services

To rebuild a specific service:
```bash
docker-compose build main-server
```

To rebuild all services:
```bash
docker-compose build
```

### Scaling Services

You can scale certain services if needed:
```bash
docker-compose up -d --scale account-server=2
```

### Debugging

1. **Access service shell:**
   ```bash
   docker-compose exec main-server sh
   ```

2. **View real-time logs:**
   ```bash
   docker-compose logs -f --tail=100 main-server
   ```

3. **Check service health:**
   ```bash
   docker-compose ps
   ```

## Stopping Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: This will delete all data)
docker-compose down -v

# Stop and remove images
docker-compose down --rmi all
```

## Troubleshooting

### Common Issues

1. **Port conflicts:** Ensure no other services are using the required ports
2. **Memory issues:** Ensure sufficient RAM is available
3. **Database connection:** Check PostgreSQL is healthy before other services start

### Health Checks

The compose file includes health checks for critical services:
- PostgreSQL: `pg_isready` command
- Valkey: `valkey-cli ping` command

### Logs Location

Application logs are stored in the `./logs` directory and are accessible from the host system.

## Production Considerations

For production deployment, consider:

1. **Security:**
   - Change default passwords
   - Use secrets management
   - Enable TLS for external connections

2. **Performance:**
   - Adjust resource limits
   - Configure proper logging levels
   - Monitor resource usage

3. **Backup:**
   - Regular database backups
   - Persistent volume backups
   - Configuration backups

## Network Architecture

All services communicate through the `agonyl-network` bridge network, enabling:
- Service discovery by name
- Isolated network communication
- Scalable architecture

## Service Dependencies

The startup order is enforced through Docker Compose dependencies:
```
Valkey → PostgreSQL → Adminer
                   ↓
                Migrate
                   ↓
              Main Server
                   ↓
            Account Server
                   ↓
            Login Server
                   ↓
             Gate Server
``` 