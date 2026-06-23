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

- [ ] **Make the board self-explanatory.** Render the coordinate inside empty cells (e.g. `a0 a1 a2 / b0 …`) so the move name is on the board. See `cmd/client/lib/display_service.go` (`PrintBoard`).
- [ ] **Fix the axis mismatch.** Rows use `a`-`c` while columns use `0`-`2` (different bases). Switch to either classic `1-9` numpad layout or 1-based `a1`-`c3` so input is unambiguous. Update `PrintBoard` and `InputService.PromptMove`.
- [ ] **Redraw in place instead of scrolling.** Clear the screen and redraw a single stable frame (board + "You are X" + whose turn) on each update, instead of appending a fresh board every poll/turn. Keep a `DEBUG`/no-TTY linear fallback.
- [ ] **Highlight what changed.** Mark/colorize the last-played cell so the opponent's move is easy to spot; highlight the winning line on a win. See `PrintBoard` / `PrintResult`.
- [ ] **Persistent status line.** Keep mark + whose turn + opponent name pinned above the board instead of printing "You are X" once at startup (`cmd/client/app/start.go`) where it scrolls away.

### Friction / controls

- [ ] **Forfeit / quit at any time.** (Also tracked in backlog above.) Input is only read on your turn in `cmd/client/app/play_game.go`, so a player waiting on an opponent cannot leave. Needs an input goroutine that listens while polling; also unblocks leaving a `waiting` game cleanly.
- [ ] **Don't pass credentials on the command line.** `make start-client USER=alice PASSWORD=alicepass` leaks the password into shell history and the process list. Prompt for the password interactively (masked) when not provided.

### Messaging

- [ ] **Friendlier connection errors.** Map `dial tcp ... connection refused` to something like "Can't reach the server — is it running on :8080?".
- [ ] **Human-readable move rejections.** `Move rejected: %v` in `play_game.go` wraps the raw API error; map known error codes to short lines ("That cell is taken", "Not your turn"). Addresses the "are server responses meaningful?" backlog item.

## Larger features (beyond polish)

- [ ] Real-time updates via WebSockets/SSE instead of 1s polling (TASK bonus).
- [ ] Ranking / leaderboard endpoint (TASK bonus).
- [ ] AFK detection + timeout-based auto-win (depends on the "last_in_game" timestamp item).
- [ ] Reconnect to an in-progress game on client restart (depends on the lookup method above).
- [ ] Matchmaking concurrency-safe without locks (SQL `SELECT FOR UPDATE`, or an in-memory queue with atomic ops/channels).
- [ ] Pass `context` to services/repository for proper cancellation and timeouts once a real DB is used.
