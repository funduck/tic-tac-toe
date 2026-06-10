package lib

import openapi "github.com/GIT_USER_ID/GIT_REPO_ID"

// Game wraps the generated openapi.GameGame with convenience methods
type Game struct {
	*openapi.GameGame
}

// GetID returns the game ID
func (g *Game) GetID() string {
	if g.GameGame == nil || g.Id == nil {
		return ""
	}
	return *g.Id
}

// GetStatus returns the game status
func (g *Game) GetStatus() string {
	if g.GameGame == nil || g.Status == nil {
		return ""
	}
	return *g.Status
}

// GetCurrentPlayerID returns the current player ID
func (g *Game) GetCurrentPlayerID() string {
	if g.GameGame == nil || g.CurrentPlayerID == nil {
		return ""
	}
	return *g.CurrentPlayerID
}

// GetUserID1 returns the first player ID
func (g *Game) GetUserID1() string {
	if g.GameGame == nil || g.UserID1 == nil {
		return ""
	}
	return *g.UserID1
}

// GetUserID2 returns the second player ID
func (g *Game) GetUserID2() string {
	if g.GameGame == nil || g.UserID2 == nil {
		return ""
	}
	return *g.UserID2
}

// GetWinnerID returns the winner ID
func (g *Game) GetWinnerID() string {
	if g.GameGame == nil || g.WinnerID == nil {
		return ""
	}
	return *g.WinnerID
}

// GetResult returns the game result
func (g *Game) GetResult() string {
	if g.GameGame == nil || g.Result == nil {
		return ""
	}
	return *g.Result
}

// GetBoard returns the game board
func (g *Game) GetBoard() [][]int32 {
	if g.GameGame == nil {
		return nil
	}
	return g.Board
}

// WrapGame wraps an openapi.GameGame into a Game
func WrapGame(apiGame *openapi.GameGame) *Game {
	return &Game{GameGame: apiGame}
}
