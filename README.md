# Gator

A CLI-based RSS feed aggregator built with Go and PostgreSQL. Manage RSS feeds, follow your favorite sources, and browse aggregated posts directly from your terminal.

## Features

- User management with account switching
- Add and follow RSS feeds
- Automatic background aggregation of posts
- Browse posts from all your followed feeds
- PostgreSQL backend with type-safe queries

## Installation

### Prerequisites

- Go 1.25.1 or higher
- PostgreSQL database

### Quick Install

```bash
go install github.com/hreshchyshynt/gator@latest
```

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build from Source

```bash
git clone https://github.com/hreshchyshynt/gator.git
cd gator
go mod download
go build -o gator
```

## Configuration

Create `~/.gatorconfig.json` with your PostgreSQL connection:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username`, `password`, and `gator` with your actual PostgreSQL credentials and database name.

**Note:** Run database migrations from `sql/schema/` before first use.

## Quick Start

```bash
# Register a user
gator register john

# Add RSS feeds
gator addfeed "Hacker News" https://news.ycombinator.com/rss
gator addfeed "Go Blog" https://go.dev/blog/feed.atom

# View all feeds
gator feeds

# Start aggregating (runs continuously, use Ctrl+C to stop)
gator agg 1m

# Browse posts (in another terminal)
gator browse 10
```

## Commands

### User Management
- `gator register <username>` - Create and login as new user
- `gator login <username>` - Switch to existing user
- `gator users` - List all users (* marks current user)
- `gator reset` - Clear all users

### Feed Management
- `gator addfeed <name> <url>` - Add feed and auto-follow
- `gator feeds` - List all feeds
- `gator follow <url>` - Follow an existing feed
- `gator following` - Show your followed feeds
- `gator unfollow <url>` - Unfollow a feed

### Aggregation & Reading
- `gator agg <interval>` - Fetch posts continuously (e.g., `1m`, `30s`, `5m`)
- `gator browse [limit]` - View recent posts (default: 2)

## Development

### Tech Stack
- Go 1.25.1
- PostgreSQL
- [sqlc](https://sqlc.dev/) for type-safe SQL queries

### Project Structure
```
gator/
├── main.go                 # Entry point
├── handlers.go             # Command handlers
├── command.go              # Command system
├── middleware.go           # Auth middleware
├── state.go                # App state
├── internal/
│   ├── config/             # Config management
│   ├── database/           # Generated queries (sqlc)
│   └── rss/                # RSS parsing
└── sql/
    ├── schema/             # Database migrations
    └── queries/            # SQL queries
```

### Database Migrations

Apply migrations in order from `sql/schema/`:
1. `001_users.sql`
2. `002_feeds.sql`
3. `003_feed_follow.sql`
4. `004_feeds_fetch_time.sql`
5. `005_posts.sql`

### Regenerate SQL Code

After modifying queries in `sql/queries/`:
```bash
sqlc generate
```

## License

Open source and available for educational purposes (Boot.dev curriculum project).
