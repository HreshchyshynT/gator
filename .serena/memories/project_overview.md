# Gator Project Overview

## Purpose
CLI-based RSS feed aggregator built with Go and PostgreSQL. Allows users to manage RSS feeds, follow sources, and browse aggregated posts from the terminal.

## Tech Stack
- **Language**: Go 1.25.1
- **Database**: PostgreSQL
- **SQL Generation**: sqlc for type-safe queries
- **Dependencies**:
  - `github.com/google/uuid` - UUID generation
  - `github.com/lib/pq` - PostgreSQL driver

## Project Structure
```
gator/
├── main.go                 # Entry point, command registration
├── command.go              # Command system (parsing, routing)
├── state.go                # Application state
├── middleware.go           # Auth middleware (middlewareLoggedIn)
├── handler_*.go            # Command handlers by domain
├── time_util.go            # Time parsing utilities
├── internal/
│   ├── config/             # Config management (~/.gatorconfig.json)
│   ├── database/           # Generated queries (sqlc)
│   └── rss/                # RSS parsing
└── sql/
    ├── schema/             # Database migrations (001-005)
    └── queries/            # SQL queries for sqlc
```

## Configuration
- Config file: `~/.gatorconfig.json`
- Contains: `db_url` (PostgreSQL connection) and `current_user_name`

## Key Commands
- User: register, login, users, reset
- Feeds: addfeed, feeds, follow, following, unfollow
- Content: agg, browse
