# Gator Project Improvements

Ideas and approaches for extending the Gator RSS aggregator project.

## 1. Browse Command Improvements

### Sorting Flags

Sorting flags allow you to view posts in different orders:
- `gator browse --sort=newest` (default behavior)
- `gator browse --sort=oldest` (see historical posts first)
- `gator browse --sort=title` (alphabetical by title)
- `gator browse --sort=feed` (group by feed source)

### Filtering Flags

Filtering flags let you narrow down what you see:
- `gator browse --limit=10` (show only 10 posts instead of 2)
- `gator browse --feed="Boot.dev Blog"` (only posts from specific feed)
- `gator browse --since="2024-01-01"` (posts after a date)
- `gator browse --unread` (if you track read status)

### Pagination

Pagination handles large result sets:
- `gator browse --page=2` (show next page of results)
- `gator browse --limit=20 --offset=40` (skip first 40, show next 20)
- Interactive: show "Press Enter for more..." and load next batch

### Combined Examples

```bash
gator browse --sort=oldest --limit=5
gator browse --feed="TechCrunch" --since="2024-12-01"
gator browse --sort=feed --page=3
```

## 2. Concurrency in Agg Command

### Current Behavior

Currently, `agg` probably processes feeds sequentially:
- Fetch feed 1 → parse → save → fetch feed 2 → parse → save, etc.

### Concurrent Approach

Concurrency would mean launching multiple feed fetches simultaneously using goroutines.

### Benefits

- If you have 10 feeds and each takes 2 seconds, sequential = 20 seconds total
- With concurrency (say 5 parallel workers), could be done in ~4-6 seconds
- More responsive when users follow many feeds

### Considerations

- How many parallel requests? (Don't want to overwhelm servers or your database)
- Database writes might need coordination (sync.Mutex or channels)
- Error handling becomes more complex (one feed failing shouldn't stop others)
- Worker pool pattern vs unbounded goroutines

**Core idea**: Yes, make fetchFeed requests in parallel, but with controlled concurrency.

## 3. Bookmarking and Liking Posts

### Data Model Approach

- Add `bookmarked` boolean column to posts table
- Or create separate `bookmarks` table linking users to posts
- Similar for `likes` - could be boolean or separate table

### User Interaction Patterns

```bash
gator bookmark <post_id>        # mark a post
gator browse --bookmarked       # view only bookmarked posts
gator unbookmark <post_id>      # remove bookmark
```

### Storage Considerations

- If multi-user system: need user_id + post_id relationship
- Single user: simple boolean flag works
- Could track timestamp of when bookmarked for sorting

### Extended Ideas

- Tags/categories for bookmarks
- Export bookmarks to file
- Sync bookmarks across devices (if you build API)

## 4. Testing HTTP API Locally

### Tools You Could Use

#### curl (command line)

```bash
curl http://localhost:8080/api/posts
curl -X POST http://localhost:8080/api/bookmark -d '{"post_id": 123}'
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/user/bookmarks
```

#### Postman

GUI tool for API testing:
- Create collections of requests
- Save authentication tokens
- Test different HTTP methods (GET, POST, PUT, DELETE)

#### HTTPie

Friendlier curl alternative:

```bash
http GET localhost:8080/api/posts
http POST localhost:8080/api/bookmark post_id=123
```

#### Browser

For simple GET requests:
- Just navigate to `http://localhost:8080/api/posts`
- Good for quick checks, not for POST/PUT/DELETE

#### Writing Tests

- Go's `net/http/httptest` package
- Create test server, make requests, verify responses
- Can test authentication, error cases, edge conditions

### Typical Workflow

1. Start your API server locally (port 8080 or whatever)
2. Use curl/Postman to make requests
3. Check responses, status codes, data format
4. Test error cases (invalid data, missing auth, etc.)

## 5. Service Manager

### What It Is

A service manager keeps long-running processes alive in the background, automatically restarting them if they crash or the system reboots.

### Common Service Managers

#### systemd (Linux)

- Most modern Linux distributions use this
- You create a `.service` file defining your program
- Commands: `systemctl start/stop/restart/status your-service`
- Auto-starts on boot if enabled

#### launchd (macOS)

- macOS's service manager
- Uses `.plist` files to configure services
- Commands: `launchctl load/unload/start/stop`

#### Windows Service (Windows)

- Built into Windows
- More complex to set up for Go programs
- Usually requires specific service wrapper code

#### supervisor (Cross-platform)

- Python-based process manager
- Simpler than systemd for some use cases
- Good for managing multiple processes

#### Docker with Restart Policies

- Containerize your app
- `docker run --restart=unless-stopped`

### Where to Read More

- **Systemd**: `man systemd.service` or search "systemd service tutorial"
- **Launchd**: Apple's Developer documentation on launchd
- **General concept**: "daemon processes" and "background services"
- **Go-specific**: Search "running Go program as systemd service"

### For Your Agg Command

The idea would be to have `agg` run continuously (not just once), checking feeds every N minutes, with the service manager ensuring it keeps running even if it crashes or the machine reboots.

---

## Additional Improvement Ideas

### TUI (Terminal User Interface)

- Use libraries like `bubbletea` or `tview` for Go
- Interactive selection of posts
- Navigate with arrow keys
- Press Enter to open in browser or view details

### Search Command

- Fuzzy searching across post titles and descriptions
- Use libraries like `fuzzysearch` or `go-fuzzywuzzy`
- Search by keywords, author, date range

### Performance Optimizations

- Database indexing on commonly queried fields
- Caching frequently accessed data
- Connection pooling for database

### Error Handling and Logging

- Better error messages for users
- Logging to file for debugging
- Metrics and monitoring
