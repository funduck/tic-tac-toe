package game

import "github.com/google/uuid"

// GameService processes game actions and delegates state persistence to GameRepo.
type GameService struct {
	repo GameRepo
}

func NewGameService(repo GameRepo) *GameService {
	return &GameService{repo: repo}
}

// CreateGame starts a new game in the waiting state, awaiting players to join.
func (s *GameService) CreateGame() (*Game, error) {
	id, err := uuid.NewV7() // V7 for lexicographically sortable IDs
	if err != nil {
		return nil, err
	}
	g := NewGame(id.String())
	if err := s.repo.Create(g); err != nil {
		return nil, err
	}
	return g, nil
}

// JoinGame adds a player to an existing waiting game. The first player to join becomes UserID1, the second becomes UserID2.
// If both players have joined, the game status changes to in_progress and UserID1 moves first.
func (s *GameService) JoinGame(gameID, userID string) (*Game, error) {
	g, err := s.repo.FindByID(gameID)
	if err != nil {
		return nil, err
	}
	if err := g.Join(userID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}

// GetLatestGameForUser retrieves the most recent game for a given user.
func (s *GameService) GetLatestGameForUser(userID string) (*Game, error) {
	return s.repo.FindLatestForUser(userID)
}

// GetGame retrieves a game by ID.
func (s *GameService) GetGame(gameID string) (*Game, error) {
	return s.repo.FindByID(gameID)
}

// MakeMove applies a move on behalf of userID and persists the updated state.
func (s *GameService) MakeMove(gameID, userID string, x, y int) (*Game, error) {
	g, err := s.repo.FindByID(gameID)
	if err != nil {
		return nil, err
	}
	if err := g.MakeMove(userID, x, y); err != nil {
		return nil, err
	}
	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}

// GiveUp concedes the game for userID, awarding the win to the other player.
func (s *GameService) GiveUp(gameID, userID string) (*Game, error) {
	g, err := s.repo.FindByID(gameID)
	if err != nil {
		return nil, err
	}
	if err := g.GiveUp(userID); err != nil {
		return nil, err
	}
	if err := s.repo.Update(g); err != nil {
		return nil, err
	}
	return g, nil
}
