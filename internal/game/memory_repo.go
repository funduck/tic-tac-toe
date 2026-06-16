package game

import (
	"context"
	"sync"
)

// MemoryRepo is a thread-safe in-memory implementation of GameRepo.
// Insertion order is tracked so FindLatestForUser returns the most recently created game.
type MemoryRepo struct {
	mu        sync.RWMutex
	gamesByID map[string]*Game
	order     []string // game IDs in creation order
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		gamesByID: make(map[string]*Game),
	}
}

// copyGame creates a deep copy of a game to prevent concurrent access issues.
func copyGame(g *Game) *Game {
	if g == nil {
		return nil
	}
	copy := *g
	return &copy
}

func (r *MemoryRepo) Create(ctx context.Context, game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	r.gamesByID[game.ID] = game
	r.order = append(r.order, game.ID)
	return nil
}

func (r *MemoryRepo) FindByID(ctx context.Context, gameID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	g, ok := r.gamesByID[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}
	return copyGame(g), nil
}

func (r *MemoryRepo) Update(ctx context.Context, game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if _, ok := r.gamesByID[game.ID]; !ok {
		return ErrGameNotFound
	}
	r.gamesByID[game.ID] = game
	return nil
}

func (r *MemoryRepo) FindLatestForUser(ctx context.Context, userID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	for i := len(r.order) - 1; i >= 0; i-- {
		g := r.gamesByID[r.order[i]]
		if g.UserID1 == userID || g.UserID2 == userID {
			return copyGame(g), nil
		}
	}
	return nil, ErrGameNotFound
}

func (r *MemoryRepo) FindGameToJoin(ctx context.Context, userID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	for i := len(r.order) - 1; i >= 0; i-- {
		g := r.gamesByID[r.order[i]]
		if g.IsJoinAllowed() && !g.Private {
			return copyGame(g), nil
		}
	}
	return nil, ErrGameNotFound
}
