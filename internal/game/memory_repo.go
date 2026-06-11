package game

import "sync"

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

func (r *MemoryRepo) Create(game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gamesByID[game.ID] = game
	r.order = append(r.order, game.ID)
	return nil
}

func (r *MemoryRepo) FindByID(gameID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	g, ok := r.gamesByID[gameID]
	if !ok {
		return nil, ErrGameNotFound
	}
	return g, nil
}

func (r *MemoryRepo) Update(game *Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.gamesByID[game.ID]; !ok {
		return ErrGameNotFound
	}
	r.gamesByID[game.ID] = game
	return nil
}

func (r *MemoryRepo) FindLatestForUser(userID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.order) - 1; i >= 0; i-- {
		g := r.gamesByID[r.order[i]]
		if g.UserID1 == userID || g.UserID2 == userID {
			return g, nil
		}
	}
	return nil, ErrGameNotFound
}

func (r *MemoryRepo) FindGameToJoin(userID string) (*Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for i := len(r.order) - 1; i >= 0; i-- {
		g := r.gamesByID[r.order[i]]
		if g.IsJoinAllowed() && !g.Private {
			return g, nil
		}
	}
	return nil, ErrGameNotFound
}
