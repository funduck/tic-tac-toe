package game

import "errors"

var ErrGameNotFound = errors.New("game not found")

// GameRepo is the persistence interface for game state.
type GameRepo interface {
	Create(game *Game) error
	FindByID(gameID string) (*Game, error)
	Update(game *Game) error
	FindLatestForUser(userID string) (*Game, error)
}
