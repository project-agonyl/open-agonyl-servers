# Agonyl MMO Game Server Project Brief

## Project Overview
Agonyl is a 2D MMO game server system built in Go with a microservices architecture. The system consists of multiple specialized servers that handle different aspects of the game.

## Core Architecture
- **Login Server**: Handles authentication and user sessions
- **Account Server**: Manages character creation, deletion, and account data
- **Main Server**: Coordinates between servers and manages zone registration
- **Gate Server**: Routes player traffic to appropriate zone servers
- **Zone Server**: Handles game world simulation and player interactions
- **Web Server**: Provides web interface for account management

## Zone Server Requirements
The zone server is the core game world simulator that needs to:

1. **Load Map Data**: Load 256x256 cell maps from binary files based on configuration
2. **TCP Server**: Implement custom protocol server for client communication
3. **World Simulation**: Handle NPC spawning, combat, movement, and events
4. **Player Management**: Process player messages and update game state
5. **Spatial Broadcasting**: Send updates to players in 9-tile radius around events

## Technical Stack
- **Language**: Go
- **Database**: PostgreSQL
- **Cache**: Redis/Valkey
- **Networking**: Custom TCP protocol
- **Data Format**: Binary files for game data
- **Architecture**: Microservices with inter-server communication

## Game Mechanics
- 2D grid-based world (256x256 cells per map)
- Real-time combat and movement
- NPC and monster spawning
- Player state management (HP, MP, skills, buffs)
- Spatial awareness and broadcasting
- Item and inventory systems

## Development Goals
- High-performance real-time simulation
- Scalable zone-based architecture
- Robust error handling and logging
- Clean separation of concerns
- Extensible design for future features 