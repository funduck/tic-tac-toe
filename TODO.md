# TODO

## Done

- [x] Unify code style across the project
- [x] Improve generated code quality by filling AGENTS.md
- [x] Return refresh token in cookie so that in browser client it is hidden from JS and not accessible to XSS attacks

## Existing backlog

- [ ] Allow forfeit for client at any time. Now it reads user's input only when it is user's turn. If there is no opponent, quit from game should just remove player and the game stays "waiting"
- [ ] Fix client bug when 2 instances try to connect with same credentials
- [ ] Save "last_in_game" timestamp for both players in the game so clients can determine when opponent is offline for a long time and the game is set to auto-win
- [ ] Add method to find/read currently playing game(s) for a user so it can reconnect to it
- [ ] Document proper client usage in README.md (Makefile is for devs only, describe DEBUG option too)
- [ ] Check TASK requirements, are they met? If not, what is missing?
- [ ] Double-check server responses, are they meaningful? Does the client display proper messages to user?

## UX polish

Proposals to make the CLI client feel like a game rather than a request log. Ordered roughly by impact/effort.

### Display layer (self-contained, low risk)

- [x] **Make the board self-explanatory.** Numpad `1`-`9` layout; empty cells show their move number. `cmd/client/lib/display_service.go`.
- [x] **Fix the axis mismatch.** Adopted classic `1-9` numpad layout; `InputService.ParseMove` accepts single keys `1`-`9`.
- [x] **Redraw in place instead of scrolling.** `RenderFrame` clears + redraws a single frame on a TTY; linear fallback under `DEBUG`/piped output. `cmd/client/lib/terminal.go`, `display_service.go`.
- [x] **Highlight what changed.** Last move (yellow) and winning line (green) highlighted. `LastMove` tracked in `state.go`, diffed in `play_game.go`.
- [x] **Persistent status line.** Header (mark · opponent · whose turn) pinned above the board in `RenderFrame`.

### Friction / controls

- [x] **Forfeit / quit at any time.** Concurrent stdin reader (`InputService.Lines`) + select loop in `play_game.go`; `q` forfeits an in-progress game via GiveUp. While still `waiting` for an opponent, `q` now cleanly leaves via the `POST /api/games/{gameID}/quit` endpoint, which cancels the game (`StatusCancelled`).
- [x] **Don't pass credentials on the command line.** `-password` is now optional; masked prompt via `golang.org/x/term` when omitted. `cmd/client/main.go`.

### Messaging

- [x] **Friendlier connection errors.** `lib.FriendlyMessage` maps network failures to "Can't reach the server — is it running on :8080?".
- [x] **Human-readable move rejections.** `lib.FriendlyMessage` maps API error codes to short lines; used in `play_game.go`/`start.go`.

## Larger features (beyond polish)

- [ ] Real-time updates via WebSockets/SSE instead of 1s polling (TASK bonus).
- [ ] Ranking / leaderboard endpoint (TASK bonus).
- [ ] AFK detection + timeout-based auto-win (depends on the "last_in_game" timestamp item).
- [ ] Reconnect to an in-progress game on client restart (depends on the lookup method above).
- [ ] Matchmaking concurrency-safe without locks (SQL `SELECT FOR UPDATE`, or an in-memory queue with atomic ops/channels).
- [ ] Pass `context` to services/repository for proper cancellation and timeouts once a real DB is used.
