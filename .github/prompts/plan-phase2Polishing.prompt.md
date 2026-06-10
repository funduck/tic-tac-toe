## Plan: Phase 2 Polishing — Logging, Error Handling, Refactoring

Improve Phase 2 implementation with structured logging, graceful error handling, meaningful OpenAPI operation IDs, and client refactoring into clean service layers.

**Steps**

### Part 1 — Logging

1. Add structured logging to `internal/game/service.go` — use `log/slog`; inject logger via constructor; log at INFO level: `CreateGame` (gameID), `JoinGame` (gameID, userID, status transition), `MakeMove` (gameID, userID, x, y, result), `GiveUp` (gameID, userID)
2. Add logging to `internal/server/handler.go` — inject logger; log requests at DEBUG level (method, path, userID from body); log errors at WARN level
3. Update `cmd/server/main.go` — create `slog.Logger` with JSON handler; pass to `GameService` and `Handler`
4. Add conditional logging middleware to `internal/server/router.go` — skip logging for `GET /games/{gameID}` (polling noise); log all other requests

### Part 2 — OpenAPI Operation IDs

5. Add operation IDs to all swaggo annotations in `internal/server/handler.go`:
   - `POST /games/create` → `@ID createGame`
   - `POST /games/{gameID}/join` → `@ID joinGame`
   - `POST /games/{gameID}/move` → `@ID makeMove`
   - `POST /games/{gameID}/giveup` → `@ID giveUpGame`
   - `GET /games/{gameID}` → `@ID getGame`
6. Run `swag init -g cmd/server/main.go` to regenerate docs
7. Run `./codegen-client.sh` to regenerate client with better method names

### Part 3 — Client Refactoring

8. Create `cmd/client/lib/display.go` — `DisplayService` with methods:
   - `PrintBoard(g *Game)`
   - `PrintResult(g *Game, userID string)`
   - `PrintStatus(status, currentPlayerID, userID string)`
   - `Clear()` (optional)
9. Create `cmd/client/lib/input.go` — `InputService` with methods:
   - `PromptMove(scanner *bufio.Scanner) (row, col int, giveUp bool, err error)` — validates input, returns parsed values
   - `Confirm(prompt string, scanner *bufio.Scanner) bool` (optional, for quit confirmation)
10. Create `cmd/client/lib/game.go` — `GameService` wrapping generated API client:
   - `CreateGame(ctx, userID) (*Game, error)`
   - `JoinGame(ctx, gameID, userID) (*Game, error)`
   - `MakeMove(ctx, gameID, userID, x, y) (*Game, error)`
   - `GiveUp(ctx, gameID, userID) (*Game, error)`
   - `GetGame(ctx, gameID) (*Game, error)`
   - `PollUntil(ctx, gameID, predicate) (*Game, error)` — handles polling with error handling
11. Create `cmd/client/lib/types.go` — define `Game` type wrapping `openapi.GameGame` with convenience methods: `GetID()`, `GetStatus()`, etc.
12. Refactor `cmd/client/main.go` — wire services together; simplified flow with error recovery in game loop (catch move errors, display, continue); no crashes on invalid moves

### Part 4 — Graceful Error Handling

13. Update client game loop in `cmd/client/main.go`:
   - Wrap all API calls in error handlers that display user-friendly messages
   - On move rejection: display error, refresh game state, continue (don't crash)
   - On polling errors: retry with backoff (max 3 retries), then display error and exit gracefully
   - On network errors: distinguish "server unreachable" vs "invalid request" messages

**Relevant files**
- `internal/game/service.go` — add logger field + logging
- `internal/server/handler.go` — add logger field + logging + operation IDs
- `internal/server/router.go` — conditional logging middleware
- `cmd/server/main.go` — create slog.Logger
- `cmd/client/lib/display.go` — new
- `cmd/client/lib/input.go` — new
- `cmd/client/lib/game.go` — new
- `cmd/client/lib/types.go` — new
- `cmd/client/main.go` — refactor to use lib services

**Dependencies**
- `log/slog` (stdlib, Go 1.21+)

**Verification**
1. `go build ./...` compiles cleanly
2. `go vet ./...` no issues
3. Server logs JSON-formatted structured logs (except polling spam)
4. Client handles invalid moves gracefully — displays error, prompts again
5. Client handles network errors gracefully — displays message, doesn't crash
6. Generated client has better method names from operation IDs
7. Two-client game still works end-to-end

**Decisions**
- Use `log/slog` for structured logging (Go 1.21+, already available in Go 1.26)
- Skip logging GET /games/{gameID} to reduce polling noise
- Client error handling: display + continue (not crash) for game logic errors; retry with backoff for network errors
- Client refactoring: lib/ package for separation of concerns
