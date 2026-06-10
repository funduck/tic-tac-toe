# Phase 1. Basic Game Workflow
## Game
Struct responsible for representing the game state, enforcing rules, and determining outcomes.

* userID1 string - uses X
* userID2 string - uses O
* board [3][3]int - 0 for empty, 1 and 2 for users
* currentPlayerID string
* status string // "waiting", "in_progress", "finished"
* result string // "win", "draw"
* winnerID string

Game methods should allow for initialization, making moves, checking for wins/draws, giving up.

## GameRepo
Interface for storing and retrieving game state.

* Create(game *Game) error
* Get(gameID string) (*Game, error)
* Update(game *Game) error
* FindLatestForUser(userID string) (*Game, error)

## GameService
Service processing game logic, managing game state, and handling player interactions.

* CreateGame(userID1, userID2 string) (*Game, error)
* MakeMove(gameID, userID string, x, y int) (*Game, error)
* GiveUp(gameID, userID string) (*Game, error)
* GetGame(gameID string) (*Game, error)

## Testing
Test Game object logic

Create mock GameRepo for testing GameService methods
Test GameService logic

# Phase 2. Client–Server Communication
## Server
HTTP server exposing endpoints for game actions.

* POST /games/create - create a new game with one player (waits for opponent)
* POST /games/{gameID}/join - join an existing game if it is waiting for an opponent
* POST /games/{gameID}/move - make a move in an active game
* POST /games/{gameID}/giveup - give up the game
* GET /games/{gameID} - get current game state. We allow anybody to read the game state for now.

### DTOs
Define requests:
* CreateGameRequest { userID string } // we are fine with providing the userID in the request body for simplicity
* JoinGameRequest { userID string, gameID string }
* MoveRequest { userID string, gameID string, x int, y int }
* GiveUpRequest { userID string, gameID string }

Annotate Game struct with JSON tags for response serialization.

### OpenAPI
Define OpenAPI spec for the above endpoints and data structures. Document request/response formats, status codes, and error messages.

Run codegenerator for client stubs `codegen_client.sh`

## Client
Command-line client to interact with the server.

* Connect to server and create/join game
* Display game board and status
* Allow user to input moves
* Display results (win/loss/draw)

