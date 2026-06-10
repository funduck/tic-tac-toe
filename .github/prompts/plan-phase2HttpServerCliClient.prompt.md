## Plan: Phase 2 — HTTP Server + CLI Client

chi router serving 5 game endpoints, swaggo annotations → `swag init` generates the spec the codegen script fetches, then a polling CLI client.

**Steps**

### Phase A — JSON tags & DTOs
1. Add JSON tags to all `Game` struct fields in `internal/game/game.go`
2. Create `internal/server/dto.go` — 4 request structs: `CreateGameRequest { GameID, UserID }`, `JoinGameRequest {  GameID, UserID }`, `MoveRequest {  GameID, UserID, X, Y }`, `GiveUpRequest {  GameID, UserID }`

### Phase B — HTTP Server
3. Create `internal/server/handler.go` — `Handler` struct holding `*GameService`; 5 methods with full swaggo annotations; error mapping: `ErrGameNotFound` → 404, `ErrGameNot*` → 409, bad JSON → 400
4. Create `internal/server/router.go` — build `chi.Router`, mount all handlers, serve swagger UI at `GET /swagger/*` and raw spec at `GET /swagger/doc.json`
5. Create `cmd/server/main.go` — top-level swag `@title`/`@host` annotations; wire `MemoryRepo → GameService → Handler`; listen on `:8080`
6. Run `swag init -g cmd/server/main.go` → generates `docs/` folder
7. Add `_ "github.com/funduck/tic-tac-toe/docs"` side-effect import to `main.go`

### Phase C — Generated Client + CLI
8. Run `./codegen_client.sh` (needs running server + Docker) → `generated/server/`
9. Create `cmd/client/main.go` — CLI with flags `--server`, `--user`, `--game`:
   - No `--game` → `POST /games/create`, print gameID, poll until opponent joins
   - With `--game` → `POST /games/{gameID}/join`. Already in game should be treated as "join and play"
   - Game loop: render board, prompt `row col`, `POST /games/{gameID}/move`; poll when waiting for opponent; display result on finish

**Relevant files**
- `internal/game/game.go` — add JSON tags
- `internal/server/dto.go` — new
- `internal/server/handler.go` — new
- `internal/server/router.go` — new
- `cmd/server/main.go` — new
- `cmd/client/main.go` — new
- `docs/` — generated
- `generated/server/` — generated

**Dependencies to add**
- `github.com/go-chi/chi/v5`
- `github.com/swaggo/swag` (+ `swag` CLI for codegen)
- `github.com/swaggo/http-swagger`

**Verification**
1. `go build ./cmd/server/` compiles
2. `go vet ./...` clean
3. `curl localhost:8080/swagger/doc.json` returns valid JSON
4. `./codegen_client.sh` generates `generated/server/` (requires Docker)
5. Two-terminal game: start server → two clients → play to win/draw

**Excluded from Phase 2**
- Auth → Phase 3
- Matchmaking queue, persistence, real-time updates, Docker → later phases
