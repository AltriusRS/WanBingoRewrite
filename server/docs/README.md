# WAN Bingo Server Documentation

A real-time bingo game platform for WAN Show streams, built with Go, PostgreSQL, and modern web technologies.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Database Schema](#database-schema)
- [Authentication](#authentication)
- [Real-time Features](#real-time-features)
- [Development](#development)
- [Deployment](#deployment)

## Overview

WAN Bingo is a live-streaming bingo game where viewers play along with WAN Show episodes. Players mark off tiles as events happen during the show, competing for the highest score and bingo achievements.

### Key Features

- **Real-time Gameplay**: Server-Sent Events for live updates
- **Guest Support**: Anonymous board generation for non-registered users
- **Discord Integration**: OAuth authentication with Discord accounts
- **Timer System**: High-accuracy countdown timers (2-second checks) with SSE notifications
- **Chat System**: Real-time messaging with threading support
- **Admin Controls**: Host tools for managing shows and tiles

## Architecture

### Tech Stack

- **Backend**: Go 1.25 with Fiber web framework
- **Frontend**: Next.js (React)
- **Database**: PostgreSQL 18 (beta) with async I/O optimizations
- **Real-time**: Server-Sent Events (SSE)
- **Authentication**: Discord OAuth2
- **Deployment**: Docker-ready

### System Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   WAN Bingo     │    │   PostgreSQL    │
│   (NextJS)      │◄──►│   Server        │◄──►│   Database      │
│                 │    │   (Go/Fiber)    │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Discord OAuth │    │   SSE Streams   │    │   Background    │
│                 │    │                 │    │   Services      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Data Flow

1. **Authentication**: Discord OAuth → Session creation
2. **Board Generation**: Random tile selection → Database storage
3. **Real-time Updates**: SSE streams for live events
4. **Timer Management**: Background monitoring → SSE notifications
5. **Chat**: WebSocket/SSE for messaging

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL 13+
- Discord Application (for OAuth)

### Setup

1. **Clone and install dependencies:**
   ```bash
   git clone <repository>
   cd wanbingo/server
   go mod download
   ```

2. **Set up PostgreSQL:**
   ```bash
   createdb wan_bingo_dev
   psql -d wan_bingo_dev < db/schema.sql
   ```

3. **Configure environment:**
   ```bash
   export DATABASE_URL="postgres://user:pass@localhost/wan_bingo_dev"
   export DISCORD_CLIENT_ID="your_client_id"
   export DISCORD_CLIENT_SECRET="your_client_secret"
   export FRONTEND_URL="http://localhost:3000"
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

5. **Verify setup:**
   ```bash
   curl http://localhost:8080/tiles | jq
   ```

## API Documentation

### REST Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/users` | List users (paginated) | Optional |
| GET | `/users/:id` | Get user profile | Optional |
| GET | `/users/me` | Get authenticated user | Required |
| GET | `/shows/latest` | Get current show | Optional |
| GET | `/shows/:id` | Get show by ID | Optional |
| GET | `/tiles` | List tiles (paginated) | Optional |
| GET | `/tiles/:id` | Get tile by ID | Optional |
| GET | `/tiles/show` | Get show tiles | Optional |
| GET | `/tiles/me` | Get user board | Required |
| GET | `/tiles/anonymous` | Generate guest board | Optional |
| GET | `/timers` | List timers (paginated) | Optional |
| GET | `/timers/:id` | Get timer | Optional |
| POST | `/timers` | Create timer | Required |
| PUT | `/timers/:id` | Update timer | Required* |
| DELETE | `/timers/:id` | Delete timer | Required* |
| POST | `/timers/:id/start` | Start timer | Required* |
| POST | `/timers/:id/stop` | Stop timer | Required* |
| POST | `/chat` | Send message | Required |
| GET | `/chat/stream` | Chat SSE stream | Optional |

*Owner only

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/auth/discord/login` | Initiate Discord OAuth |
| GET | `/auth/discord/callback` | Handle OAuth callback |
| POST | `/auth/discord/logout` | Clear session |

### Real-time Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/chat/stream` | Chat messages & events | Optional |
| GET | `/host/stream` | Administrative events | Required* |

*Host permissions required

## Database Schema

### Core Tables

- **`players`** - User accounts (Discord OAuth)
- **`sessions`** - Authentication tokens
- **`shows`** - WAN Show episodes
- **`tiles`** - Bingo squares
- **`boards`** - Player bingo cards
- **`messages`** - Chat messages
- **`timers`** - Countdown timers

### Relationships

```
players (1) ──── (N) sessions
players (1) ──── (N) boards
players (1) ──── (N) messages
players (1) ──── (N) timers

shows (1) ──── (N) boards
shows (1) ──── (N) messages
shows (1) ──── (N) timers
shows (1) ──── (N) show_tiles

tiles (1) ──── (N) show_tiles
tiles (1) ──── (N) tile_confirmations
```

### Key Constraints

- One board per player per show
- Unique Discord ID per player
- Soft deletes on all tables
- Automatic timestamp updates

## Authentication

### Discord OAuth2 Flow

1. User clicks "Login with Discord"
2. Server redirects to Discord authorization
3. Discord redirects back with authorization code
4. Server exchanges code for access token
5. Server fetches user info and creates/updates player
6. Server creates session and sets cookie
7. User redirected to frontend

### Session Management

- **Duration**: 24 hours
- **Storage**: HttpOnly, Secure cookies
- **Validation**: Database lookup on each request
- **Cleanup**: Automatic expiration handling

### Permissions

```json
{
  "canChat": true,
  "canHost": false,
  "canModerate": false,
  "canManageTimers": false
}
```

## Real-time Features

### Server-Sent Events (SSE)

The application uses SSE for real-time communication:

- **Connection Types**: Chat streams, host streams
- **Event Format**: JSON with opcode and data
- **Reconnection**: Automatic client-side handling
- **Authentication**: Permission-based event filtering

### Event Types

#### Chat Events
- `hub.connected` - Connection established
- `hub.authenticated` - Permission information
- `hub.connections.count` - User count updates
- `chat.message` - New messages
- `chat.players` - Participant information

#### Show Events
- `whenplane.aggregate` - Show status updates

#### Timer Events
- `timer.expired` - Timer completion (host stream)

### Background Services

- **Timer Monitor**: Checks for expired timers every 2 seconds
- **Session Cleanup**: Removes expired authentication tokens
- **Event Broadcasting**: Distributes events to connected clients

## Development

### Project Structure

```
server/
├── db/                    # Database layer
│   ├── models/           # Data models
│   ├── board.go          # Board operations
│   ├── player.go         # Player operations
│   ├── session.go        # Session management
│   ├── show.go           # Show operations
│   ├── tile.go           # Tile operations
│   ├── timer.go          # Timer operations
│   └── schema.sql        # Database schema
├── handlers/             # HTTP handlers
│   ├── auth/            # Authentication
│   ├── chat/            # Chat endpoints
│   ├── show/            # Show endpoints
│   ├── tiles/           # Tile endpoints
│   ├── timers/          # Timer endpoints
│   └── users/           # User endpoints
├── middleware/           # HTTP middleware
│   ├── discord.go       # Auth middleware
│   └── logger.go        # Request logging
├── sse/                 # Server-Sent Events
│   ├── client.go        # SSE client handling
│   ├── hub.go           # Event broadcasting
│   └── chat.go          # Chat-specific SSE
├── timers/              # Timer background services
│   └── monitor.go       # Timer expiration monitoring
├── utils/               # Utilities
│   ├── utils.go         # Common functions
│   └── log.go           # Logging helpers
├── docs/                # Documentation
├── main.go              # Application entry point
└── go.mod               # Go dependencies
```

### Development Commands

```bash
# Run tests
go test ./...

# Build for development
go build -o server-dev

# Run with hot reload
air  # Install with: go install github.com/cosmtrek/air@latest

# Database migrations
psql -d wan_bingo_dev < db/schema.sql

# Check code quality
gofmt -l .
go vet ./...
```

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://user:pass@localhost/wan_bingo_dev

# Discord OAuth
DISCORD_CLIENT_ID=your_client_id
DISCORD_CLIENT_SECRET=your_client_secret
DISCORD_REDIRECT_URI=http://localhost:8080/auth/discord/callback

# Application
PORT=8080
FRONTEND_URL=http://localhost:3000

# Development
DEBUG=true
LOG_LEVEL=debug
```

## Deployment

### Docker Deployment

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/db/schema.sql .
CMD ["./server"]
```

### Production Checklist

- [ ] Set `SECURE=true` on session cookies
- [ ] Configure HTTPS/TLS
- [ ] Set up database backups
- [ ] Configure monitoring/alerting
- [ ] Set up log aggregation
- [ ] Configure rate limiting
- [ ] Set up health checks
- [ ] Configure CORS properly

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Database connectivity
curl http://localhost:8080/health/db

# SSE connectivity
curl http://localhost:8080/chat/stream
```

## Contributing

### Code Style

- Follow Go conventions (`gofmt`, `go vet`)
- Use meaningful variable names
- Add comments for exported functions
- Write tests for new features
- Update documentation

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Update documentation
6. Submit pull request

### Commit Messages

```
feat: add timer expiration SSE events
fix: handle edge case in board generation
docs: update API documentation
refactor: simplify authentication middleware
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Documentation**: This docs folder
- **Community**: Discord server (link in README)

---

## Additional Resources

- [API Reference](api.md) - Complete endpoint documentation
- [Database Models](models.md) - Schema and relationships
- [Authentication Guide](auth.md) - OAuth setup and session management
- [Real-time Events](realtime.md) - SSE event documentation
- [Database Guide](database.md) - Schema, migrations, and maintenance</content>
</xai:function_call">Main README Documentation