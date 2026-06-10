## Plan: Phase 1 — Core Game Domain

Implement the `Game` domain object, `GameRepo` interface + in-memory implementation, and `GameService`, then cover everything with unit tests. No HTTP, no persistence — pure logic.

**Steps**

### Phase A — Domain Model

1. Create `internal/game/game.go` — define `Game` struct (with `id string` field added to the roadmap list), constants for statuses/results, and all methods: `NewGame`, `MakeMove`, `GiveUp`, `checkWin`, `checkDraw`
2. Create `internal/game/game_test.go` — table-driven tests for all `Game` methods: valid moves, win detection (rows/cols/diags), draw, invalid moves (occupied, out-of-bounds, wrong player, game already finished)

### Phase B — Repository

3. Create `internal/game/repo.go` — define `GameRepo` interface with `Create`, `Get`, `Update`, `FindLatestForUser`
4. Create `internal/game/memory_repo.go` — concrete in-memory implementation using a `sync.RWMutex`-guarded map; `FindLatestForUser` scans all games and returns the most recent by creation order

### Phase C — Service

5. Create `internal/game/service.go` — `GameService` struct with `GameRepo` dependency; implement `CreateGame`, `MakeMove`, `GiveUp`, `GetGame`; `CreateGame` generates a UUID7 for game ID
6. Create `internal/game/service_test.go` — use a hand-written mock `GameRepo` (implements the interface); test all service methods including error paths (game not found, wrong player, invalid move forwarded to `Game`)

**Relevant files**
- `go.mod` — module is `github.com/funduck/tic-tac-toe`, Go 1.26
- `internal/game/game.go` — new
- `internal/game/game_test.go` — new
- `internal/game/repo.go` — new
- `internal/game/memory_repo.go` — new
- `internal/game/service.go` — new
- `internal/game/service_test.go` — new

**Verification**
1. `go test ./internal/game/...` passes with no failures
2. `go vet ./...` reports no issues
3. Test cases explicitly cover: all 8 win conditions, draw, each invalid-move error code, `GiveUp` by both players, `GiveUp` on finished game (rejected)

**Decisions**
- `id` field added to `Game` struct; UUID7 generated in `GameService.CreateGame`
- In-memory repo included as a concrete impl (will be reused in Phase 2 server)
- github.com/google/uuid for uuid v7
- Auth, matchmaking, HTTP deferred to later phases
- Errors: use `errors.New` sentinels for now (e.g. `ErrNotYourTurn`, `ErrCellOccupied`) — simple and effective for core logic; can evolve to typed errors if needed in later phases
