package game

import "errors"

const (
	StatusWaiting    = "waiting"
	StatusInProgress = "in_progress"
	StatusFinished   = "finished"

	ResultWin  = "win"
	ResultDraw = "draw"
)

var (
	ErrNotYourTurn    = errors.New("not your turn")
	ErrCellOccupied   = errors.New("cell is occupied")
	ErrOutOfBounds    = errors.New("move out of bounds")
	ErrGameNotWaiting = errors.New("game is not waiting for players")
	ErrGameNotActive  = errors.New("game is not in progress")
	ErrNotInGame      = errors.New("player not in this game")
)

// Game represents the state of a single tic-tac-toe match.
// UserID1 plays as X (mark 1), UserID2 plays as O (mark 2).
type Game struct {
	ID              string    `json:"id"`
	UserID1         string    `json:"userID1"`
	UserID2         string    `json:"userID2"`
	Board           [3][3]int `json:"board"`
	CurrentPlayerID string    `json:"currentPlayerID"`
	Status          string    `json:"status"`
	Result          string    `json:"result"`
	WinnerID        string    `json:"winnerID"`
	Private         bool      `json:"private"` // optional field to indicate if the game is private or public
}

// NewGame creates a new game in the waiting state.
func NewGame(id string) *Game {
	return &Game{
		ID:     id,
		Status: StatusWaiting,
	}
}

// NewGameInProgress creates a game that is immediately in progress.
// UserID1 moves first.
func NewGameInProgress(id, userID1, userID2 string) *Game {
	return &Game{
		ID:              id,
		UserID1:         userID1,
		UserID2:         userID2,
		CurrentPlayerID: userID1,
		Status:          StatusInProgress,
	}
}

func (g *Game) IsJoinAllowed() bool {
	if g.Status != StatusWaiting {
		return false
	}
	if g.UserID1 != "" && g.UserID2 != "" {
		return false
	}
	return true
}

// Join adds a player to the game. The first player to join becomes UserID1, the second becomes UserID2.
// If both players have joined, the game status changes to in_progress and UserID1 moves first.
func (g *Game) Join(userID string) error {
	if userID == g.UserID1 || userID == g.UserID2 {
		return nil // already joined, no change
	}
	if g.Status != StatusWaiting {
		return ErrGameNotWaiting
	}
	if g.UserID1 == "" {
		g.UserID1 = userID
	} else {
		g.UserID2 = userID
	}
	if g.UserID1 != "" && g.UserID2 != "" {
		g.Status = StatusInProgress
		g.CurrentPlayerID = g.UserID1
	}
	return nil
}

// MakeMove places the mark for userID at (x, y).
// Returns an error for any invalid move without modifying state.
func (g *Game) MakeMove(userID string, x, y int) error {
	// do not allow anyone except players in the game
	if g.playerMark(userID) == 0 {
		return ErrNotInGame
	}

	if g.Status != StatusInProgress {
		return ErrGameNotActive
	}
	if g.CurrentPlayerID != userID {
		return ErrNotYourTurn
	}
	if x < 0 || x > 2 || y < 0 || y > 2 {
		return ErrOutOfBounds
	}
	if g.Board[x][y] != 0 {
		return ErrCellOccupied
	}

	mark := g.playerMark(userID)
	g.Board[x][y] = mark

	if g.checkWin(mark) {
		g.Status = StatusFinished
		g.Result = ResultWin
		g.WinnerID = userID
		return nil
	}
	if g.checkDraw() {
		g.Status = StatusFinished
		g.Result = ResultDraw
		return nil
	}

	g.CurrentPlayerID = g.otherPlayer(userID)
	return nil
}

// GiveUp ends the game, awarding victory to the other player.
func (g *Game) GiveUp(userID string) error {
	// do not allow anyone except players in the game
	if g.playerMark(userID) == 0 {
		return ErrNotInGame
	}

	if g.Status != StatusInProgress {
		return ErrGameNotActive
	}
	g.Status = StatusFinished
	g.Result = ResultWin
	g.WinnerID = g.otherPlayer(userID)
	return nil
}

func (g *Game) playerMark(userID string) int {
	switch userID {
	case g.UserID1:
		return 1
	case g.UserID2:
		return 2
	default:
		return 0
	}
}

func (g *Game) otherPlayer(userID string) string {
	if userID == g.UserID1 {
		return g.UserID2
	}
	return g.UserID1
}

func (g *Game) checkWin(mark int) bool {
	b := g.Board
	for i := 0; i < 3; i++ {
		if b[i][0] == mark && b[i][1] == mark && b[i][2] == mark {
			return true
		}
		if b[0][i] == mark && b[1][i] == mark && b[2][i] == mark {
			return true
		}
	}
	if b[0][0] == mark && b[1][1] == mark && b[2][2] == mark {
		return true
	}
	if b[0][2] == mark && b[1][1] == mark && b[2][0] == mark {
		return true
	}
	return false
}

func (g *Game) checkDraw() bool {
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			if g.Board[r][c] == 0 {
				return false
			}
		}
	}
	return true
}
