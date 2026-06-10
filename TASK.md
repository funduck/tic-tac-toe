# Go Developer

## Technical Assessment

#### Multiplayer Tic Tac Toe

### Overview

Build a multiplayer Tic Tac Toe game server in Go. This challenge evaluates your ability to write clean,
idiomatic Go code, design clear APIs, handle client–server communication, and write meaningful tests. You
are free to use any supporting technology — just keep Go as the primary language.

### Minimum Requirements

#### Game Flow

- Start the server, then run the client app in two separate terminals.
- The two clients automatically connect and start a game.
- Players alternate turns until one wins or the game ends in a draw.
- The final result (win / loss / draw) is clearly displayed to both clients.

#### Game Logic

- Enforce the standard 3×3 Tic Tac Toe grid.
- Alternate turns strictly between the two players.
- Reject invalid moves (occupied cell, out-of-turn, out-of-bounds).
- Correctly detect all win conditions and draw state.

#### Client–Server Communication

- Clients must be able to send moves and receive game-state updates.
- The protocol must be clearly defined — document it in your README.
- The server must support multiple simultaneous games (see Bonus).

### Authentication (Required)

Before joining the matchmaking queue or making any move, every client must authenticate with the server.

- Implement a register and login flow (endpoint, command, or handshake — your choice of protocol).
- On successful login the server issues a token (JWT, session, API key — justify your choice).
- All subsequent operations — join queue, make move, view leaderboard — must carry a valid token.
- Unauthenticated or invalid requests must be rejected with a clear error.
- Document your auth design decisions in the README.


```
n There is no single correct mechanism. Pick what fits your architecture and be ready to explain the tradeoffs.
```
### Bonus Objectives

```
These are optional but will strengthen your submission:
Matchmaking — A queue-based system that automatically pairs two waiting players.
Ranking / Leaderboard — Track wins, losses, draws — expose a leaderboard endpoint or command.
Persistence — Store game state and rankings in a DB (SQLite, Postgres, etc.). Support server restarts.
Real-Time Updates — Use WebSockets or SSE instead of polling for live game updates.
Deployment — Provide a Dockerfile or docker-compose.yml for easy setup.
```
### Evaluation Criteria

Correctness Rules enforced, end-to-end game flow works.

Authentication Auth is complete, tokens are validated, errors are clear.

API Design & UX Protocol is clean and easy to use from the client side.

Code Quality Idiomatic, well-structured, maintainable Go.

Testing Meaningful unit tests covering core game and auth logic.

Bonus Any of the optional objectives implemented well.

### Submission

```
Please submit your solution as:
```
- A link to a public GitHub or GitLab repository, or
- A .zip archive with the full project (exclude compiled binaries).
Include a README.md covering: how to run the project, your protocol design, auth mechanism choice, and
any architectural decisions worth explaining.

```
Good luck, and have fun building! n
```

