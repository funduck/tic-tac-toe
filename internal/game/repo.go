package game

import (
	"context"
	"errors"
)

var (
	ErrGameNotFound    = errors.New("game not found")
	ErrVersionConflict = errors.New("game was modified concurrently")
)

// GameRepo is the persistence interface for game state.
type GameRepo interface {
	Create(ctx context.Context, game *Game) error
	FindByID(ctx context.Context, gameID string) (*Game, error)
	Update(ctx context.Context, game *Game) error
	FindLatestForUser(ctx context.Context, userID string) (*Game, error)
	FindGameToJoin(ctx context.Context, userID string) (*Game, error)
}
