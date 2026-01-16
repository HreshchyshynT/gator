# Suggested Commands

## Build and Run
```bash
# Build the project
go build -o gator

# Run directly without building
go run .

# Install globally
go install
```

## Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test file
go test -v command_test.go
```

## Code Quality
```bash
# Format code
go fmt ./...

# Lint (if golangci-lint installed)
golangci-lint run

# Vet for common mistakes
go vet ./...
```

## Database/SQL
```bash
# Regenerate sqlc queries after modifying sql/queries/
sqlc generate

# Apply migrations (manual - run in psql)
psql -d gator -f sql/schema/001_users.sql
psql -d gator -f sql/schema/002_feeds.sql
# etc.
```

## Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

## Git
```bash
git status
git diff
git add .
git commit -m "message"
```

## Application Usage
```bash
# Register user
./gator register <username>

# Add feed
./gator addfeed "<name>" <url>

# Browse posts
./gator browse [limit] [--sort=newest|oldest|title|feed] [--feed=<name>] [--since=<date>]

# Start aggregator
./gator agg <interval>  # e.g., 1m, 30s
```
