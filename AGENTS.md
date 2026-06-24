# Code style
Keep idiomatic Go formatting

# Tests
When reasonable avoid individual test functions in favor of table-driven tests.

# Error variables
Group multiple package-level error vars into a single `var (...)` block.
```go
// correct
var (
    ErrFoo = errors.New("foo")
    ErrBar = errors.New("bar")
)

// wrong
var ErrFoo = errors.New("foo")
var ErrBar = errors.New("bar")
```

# Log messages and keys
- Log messages are lowercase (`"user signed up"`, not `"User signed up"`)
- Structured log keys use camelCase (`"userID"`, not `"user_id"`)

# Constructor parameter names
Constructor parameter names must match the struct field they initialise.
```go
// correct
func NewFooHandler(svc FooService, logger *slog.Logger) *FooHandler {
    return &FooHandler{svc: svc, logger: logger}
}

// wrong — parameter called "fooService" but field is "svc"
func NewFooHandler(fooService FooService, ...) *FooHandler {
    return &FooHandler{svc: fooService, ...}
}
```

# Receiver names
- Services use `s`
- Handlers use `h`
- Repos use `r`
```go
func (s *GameService) ...
func (h *GameHandler) ...
func (r *MemoryRepo) ...
```

# Router parameter names
Router functions receive a handler as `h`, not an abbreviation of the type name.
```go
// correct
func UserRouter(r chi.Router, h *UserHandler, ...)

// wrong
func UserRouter(r chi.Router, uh *UserHandler, ...)
```

# Go doc comments
Doc comments must start with the name of the function or type they describe.
```go
// ValidateToken parses and validates the access token...
func (s *AccessTokenService) ValidateToken(...) ...

// wrong — does not start with function name
// Parses and validates the access token...
func (s *AccessTokenService) ValidateToken(...) ...
```

# Swagger annotations
- `@Success` status codes must match the actual `http.Status*` constant used in the handler
- `@ID` alignment uses a single tab, consistent with all other annotation tags

# JSON tags
Use idiomatic snake_case for json: `json:"access_token"`, `json:"user_id"`

# Tools
Look into @Makefile for available commands like lint, tests, etc.

## Dev Notes
- OpenAPI client (`cmd/client/generated/server`) is generated from the *live* server's `/swagger/doc.json`, not the source: `make swag-init` (annotations → docs/swagger.json) then start the server + `make codegen-client` (Docker required).
- Generated client maps Go `time.Time` json fields to `*string` (swagger 2.0 has no time type).
- `memory_repo.copyGame` is a shallow struct copy; pointer fields (e.g. `*time.Time`) are safe only because callers reassign the pointer, never mutate through it. `Update` replaces the whole map entry.
- `cmd/client` is a separate Go module (`github.com/funduck/tic-tac-toe/client`); run its tests from that dir — root `go test ./...` won't cover it.
- AFK auto-win is enforced lazily inside `GetGame` (active player's poll calls `Touch` then `ForfeitIfOpponentAFK`), not a background reaper — the AFK player only loses once the opponent polls. Timeout via `AFK_TIMEOUT` (default 30s).