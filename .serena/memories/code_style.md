# Code Style and Conventions

## Naming
- **Functions**: camelCase (e.g., `handleBrowse`, `parseArguments`)
- **Types**: PascalCase (e.g., `Command`, `browseOptions`)
- **Constants**: camelCase (e.g., `sortNewest`, `limitArg`)
- **Packages**: lowercase single words

## Handler Pattern
- Handlers named `handle<Action>` (e.g., `handleBrowse`, `handleLogin`)
- Handler signature: `func(s *State, command Command) error`
- Protected handlers use middleware wrapper: `middlewareLoggedIn(handler)`
- Protected handler signature: `func(s *State, command Command, user database.User) error`

## File Organization
- One handler file per domain: `handler_browse.go`, `handler_user.go`, etc.
- Internal packages for reusable code: config, database, rss
- SQL queries in `sql/queries/`, schemas in `sql/schema/`

## Error Handling
- Return errors with context: `fmt.Errorf("Error get posts from db: %v", err)`
- Use `log.Fatalf` or `log.Fatalln` for fatal errors in main
- Check errors immediately after function calls

## Imports
- Standard library first, then external, then internal
- Use import aliases when needed: `cfg "github.com/hreshchyshynt/gator/internal/config"`

## SQL/Database
- Use sqlc for type-safe queries
- Generated code in `internal/database/`
- Use `sql.NullString`, `sql.NullTime` for optional parameters
- Context passed to all database calls

## Command Arguments
- Support named arguments: `--name=value`
- Support positional arguments for backward compatibility
- Parse in dedicated `parseArguments` functions
