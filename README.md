# Multiplayer Tic Tac Toe

A multiplayer Tic Tac Toe game server and client implementation in Go, featuring JWT-based authentication, automatic matchmaking, and a RESTful API.

**Contents:**
- [Quick Start](#quick-start)
- [Protocol Design](#protocol-design)
- [Authentication](#authentication)
- [Development](#development)
- [Future Improvements](#future-improvements)

## Quick Start

### Prerequisites

- Go 1.26.2 or higher
- `make` (optional, for convenience commands)

### Installation

Clone the repository and install dependencies:

```bash
git clone <repository-url>
cd tic-tac-toe
go mod download
```

### Running the Server

Start the server on `localhost:8080`:

```bash
# Using make
make start-server

# Or directly
go run cmd/server/main.go
```

The server will start and display:
```
Server listening on :8080
```

You can access the Swagger UI documentation at: `http://localhost:8080/swagger/index.html`

### Running the Client

Open **two separate terminals** and run a client in each:

**Terminal 1 (Player 1):**
```bash
go run cmd/client/main.go --user alice --password alicepass
```

**Terminal 2 (Player 2):**
```bash
go run cmd/client/main.go --user bob --password bobpass
```

The clients will:
1. Automatically register (on first run) or login
2. Create or join a game via matchmaking
3. Take turns making moves
4. Display the game board after each move
5. Show the final result (win/loss/draw)

### Client Options

```bash
--server string      Server base URL (default "http://localhost:8080")
--user string        Your user ID (required)
--password string    Your password (required)
--game string        Game ID to join a specific game (optional)
--private            Create a private game (default false)
```

### Running Tests

Run all tests with verbose output:

```bash
make tests

# Or directly
go test -v ./...
go vet ./...
```

Run tests with race detection:

```bash
go test -race ./...
```

---

## Protocol Design

### Overview

The system uses a **REST API over HTTP** with **polling-based state updates**. Clients communicate with the server via HTTP POST/GET requests, and poll the game state periodically to detect changes (opponent moves, game completion, etc.).

### Game Flow

```mermaid
sequenceDiagram

ClientA ->> Server: POST /api/users/signup (alice)
Server -->> ClientA: 200 OK + tokens    

ClientA ->> Server: POST /api/games (create game)
Server -->> ClientA: 201 Created + gameID

loop Wait for opponent
    ClientA ->> Server: GET /api/games/{id} (polling for opponent)
    Server -->> ClientA: 200 OK (waiting)   
end

ClientB ->> Server: POST /api/users/signup (bob)
Server -->> ClientB: 200 OK + tokens
ClientB ->> Server: POST /api/games/join (matchmaking)
Server -->> ClientB: 200 OK + gameID

loop Wait for turn
    ClientB ->> Server: GET /api/games/{id} (polling for turn)
    Server -->> ClientB: 200 OK (your turn)
end

ClientA ->> Server: POST /api/games/{id}/move (make move)
Server -->> ClientA: 200 OK


Note over ClientA, ClientB: game continues until the result, for example Alice wins

ClientA ->> Server: POST /api/games/{id}/move (winning move)
Server -->> ClientA: 200 OK (finished, winner: alice)
```

### Authentication Flow

#### Signup
```mermaid
sequenceDiagram
Client ->> Server: POST /api/users/signup (userID, password)
Server -->> Client: 201 Created + {access_token, refresh_token}

Note over Client, Server: Client stores tokens securely (in memory for CLI, localStorage for web)
```

#### Login
```mermaid
sequenceDiagram

Client ->> Server: POST /api/users/login (userID, password)
Server -->> Client: 200 OK + {access_token, refresh_token}

Note over Client, Server: Same as signup, but for existing users
```

#### Using Access Token
```mermaid
sequenceDiagram

Client ->> Server: POST /api/games (Authorization: Bearer access_token)
Server -->> Client: 201 Created + gameID    
Client ->> Server: GET /api/games/{id} (Authorization: Bearer access_token)
Server -->> Client: 200 OK + game state

Note over Client, Server: When token expires, client receives 401 Unauthorized and must use refresh token

Client ->> Server: POST /api/games/{id}/move (Authorization: Bearer access_token, move data)
Server -->> Client: 401 Unauthorized (if token expired)
Client ->> Server: POST /api/users/refresh-token (userID, refreshToken)
Server -->> Client: 200 OK + {new_access_token, new_refresh_token}
Client ->> Server: POST /api/games/{id}/move (repeats original request with new access token)

Note over Client, Server: This flow ensures seamless user experience even with short-lived access tokens

```

### API Endpoints
**Run server** and open [OpenAPI docs](http://localhost:8080/swagger/index.html) for detailed endpoint information.

### Polling Mechanism

Clients use a **polling strategy** to monitor game state changes:

- **Wait for opponent:** Poll every 1 second until `status` changes from `"waiting"` to `"in_progress"`
- **Wait for turn:** Poll every 1 second until `currentPlayerID` matches your user ID
- **Check game end:** Poll every 1 second until `status` becomes `"finished"`

This approach is simple and stateless but introduces latency (up to 1 second for updates). Future improvements could include WebSockets or Server-Sent Events for real-time updates.

### Error Handling

All errors return HTTP status codes with a JSON body:

```json
{
  "error": "error message description"
}
```

Common status codes:
- `400 Bad Request` - Invalid input, occupied cell, out-of-turn move, etc.
- `401 Unauthorized` - Missing or invalid authentication token
- `404 Not Found` - Game or user not found
- `409 Conflict` - User already exists (signup)
- `500 Internal Server Error` - Unexpected server error

---

## Authentication

### Mechanism: JWT (JSON Web Tokens)

The server uses **JWT tokens** for stateless authentication. This design choice offers several advantages:

**Why JWT?**
- ✅ **Stateless:** No server-side session storage required (scales horizontally)
- ✅ **Standard:** Well-established, widely supported protocol
- ✅ **Self-contained:** Tokens carry user identity (claims), reducing database lookups
- ✅ **Secure:** HMAC-SHA256 signature prevents tampering
- ✅ **Expiration:** Built-in time-based validity

**Tradeoffs:**
- ❌ Cannot revoke tokens before expiration (mitigation: short-lived access tokens + refresh tokens stored in the database)
- ❌ Token size larger than session IDs (acceptable overhead for HTTP headers)

### Token Lifecycle

The system uses a **dual-token approach**:

1. **Access Token**
   - Lifetime: **10 seconds**
   - Used for authenticating game API requests
   - Short-lived to minimize risk if leaked
   - Claims: `user_id`, `exp`, `iat`, `nbf`, `iss`

2. **Refresh Token**
   - Lifetime: **7 days**
   - Used to obtain new access tokens without re-login
   - Stored hashed on server side for security
   - Cannot be used for game operations

### Security Considerations

- **Password storage:** Passwords are hashed using **bcrypt** (cost factor 10) before storage
- **Minimum password length:** 6 characters
- **Token signing:** HMAC-SHA256 with a secret key (configurable)
- **Refresh token storage:** Refresh tokens are hashed (SHA256) before storage to prevent misuse if database is compromised
- **Token validation:** Every protected endpoint validates the Bearer token signature, expiration, and issuer

---

## Development

### Hot Reload (Server)

Use `air` for automatic server reloads during development:

```bash
# Install air
go install github.com/air-verse/air@latest

# Run with hot reload
make serve  # or just 'air'
```

### Regenerate OpenAPI Client

If you modify the API endpoints:

```bash
# Update Swagger docs
make swag-init

# Regenerate client stubs
make codegen-client
```

--- 

## Future Improvements

Additionally to bonus features mentioned in the [task](./TASK.md), here are some future improvements:
* If `-game` parameter is not passed to client, it looks up for the last played unfinished game and joins it instead of creating a new one. This allows for easier testing and quicker game start without needing to copy-paste game IDs.
* Store in memory the last timestamp when user *touched* the game and provide this info to the opponent. It allows to detect if opponent is AFK and maybe implement some timeout-based game end in the future.