package game

import "errors"

var ErrGameNotFound = errors.New("game not found")

// GameRepo is the persistence interface for game state.
type GameRepo interface {
	Create(game *Game) error
	FindByID(gameID string) (*Game, error)
	Update(game *Game) error
	FindLatestForUser(userID string) (*Game, error)
	FindGameToJoin(userID string) (*Game, error) // optional method to find a waiting game for a user to join
}
