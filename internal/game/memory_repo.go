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

func (r *MemoryRepo) Create(ctx context.Context, game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}
	r.gamesByID[game.ID] = game.Clone()
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
	return g.Clone(), nil
}

func (r *MemoryRepo) Update(ctx context.Context, game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}
	stored, ok := r.gamesByID[game.ID]
	if !ok {
		return ErrGameNotFound
	}
	if stored.Version != game.Version {
		return ErrVersionConflict
	}
	game.Version++
	r.gamesByID[game.ID] = game.Clone()
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
			return g.Clone(), nil
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
			return g.Clone(), nil
		}
	}
	return nil, ErrGameNotFound
}
