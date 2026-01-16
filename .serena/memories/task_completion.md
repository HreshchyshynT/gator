# Task Completion Checklist

## Before Committing Code

1. **Format code**
   ```bash
   go fmt ./...
   ```

2. **Run vet**
   ```bash
   go vet ./...
   ```

3. **Run tests**
   ```bash
   go test ./...
   ```

4. **Build successfully**
   ```bash
   go build -o gator
   ```

5. **If SQL queries modified**
   ```bash
   sqlc generate
   ```

## Code Review Checklist
- [ ] Error handling present for all error-returning functions
- [ ] Context passed to database calls
- [ ] New handlers follow `handle<Action>` naming
- [ ] Protected routes use `middlewareLoggedIn` wrapper
- [ ] SQL changes regenerated with `sqlc generate`
- [ ] No hardcoded values (use constants)

## Testing New Features
1. Build: `go build -o gator`
2. Test command manually: `./gator <command>`
3. Run automated tests: `go test ./...`
