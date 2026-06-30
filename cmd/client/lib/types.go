package lib

import (
	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

// Game wraps the generated openapi.GameGame with convenience methods
type Game struct {
	*openapi.GameGame
}

func (g *Game) GetOpponentID() string {
	if g.GameGame == nil {
		return ""
	}
	if g.UserId1 != nil && *g.UserId1 != UserID {
		return *g.UserId1
	}
	if g.UserId2 != nil && *g.UserId2 != UserID {
		return *g.UserId2
	}
	return ""
}

// GetOpponentLastSeen returns the opponent's last-seen timestamp (RFC3339), or nil.
func (g *Game) GetOpponentLastSeen() *string {
	if g.GameGame == nil {
		return nil
	}
	if g.UserId1 != nil && *g.UserId1 != UserID {
		return g.UserId1LastSeen
	}
	if g.UserId2 != nil && *g.UserId2 != UserID {
		return g.UserId2LastSeen
	}
	return nil
}

// GetID returns the game ID
func (g *Game) GetID() string {
	if g.GameGame == nil || g.Id == nil {
		return ""
	}
	return *g.Id
}

// GetStatus returns the game status
func (g *Game) GetStatus() openapi.GameGameStatus {
	if g.GameGame == nil || g.Status == nil {
		return ""
	}
	return *g.Status
}

// GetCurrentPlayerID returns the current player ID
func (g *Game) GetCurrentPlayerID() string {
	if g.GameGame == nil || g.CurrentPlayerId == nil {
		return ""
	}
	return *g.CurrentPlayerId
}

// GetUserID1 returns the first player ID
func (g *Game) GetUserID1() string {
	if g.GameGame == nil || g.UserId1 == nil {
		return ""
	}
	return *g.UserId1
}

// GetUserID2 returns the second player ID
func (g *Game) GetUserID2() string {
	if g.GameGame == nil || g.UserId2 == nil {
		return ""
	}
	return *g.UserId2
}

// GetWinnerID returns the winner ID
func (g *Game) GetWinnerID() string {
	if g.GameGame == nil || g.WinnerId == nil {
		return ""
	}
	return *g.WinnerId
}

// GetResult returns the game result
func (g *Game) GetResult() openapi.GameGameResult {
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
