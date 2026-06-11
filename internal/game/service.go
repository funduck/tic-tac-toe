package game

import (
	"log/slog"

	"github.com/google/uuid"
)

// GameService processes game actions and delegates state persistence to GameRepo.
type GameService struct {
	repo   GameRepo
	logger *slog.Logger
}

func NewGameService(repo GameRepo, logger *slog.Logger) *GameService {
	return &GameService{repo: repo, logger: logger}
}

func (s *GameService) createGame(private bool) (*Game, error) {
	id, err := uuid.NewV7() // V7 for lexicographically sortable IDs
	if err != nil {
		return nil, err
	}
	g := NewGame(id.String())
	if private {
		g.Private = true
	}
	if err := s.repo.Create(g); err != nil {
		return nil, err
	}
	s.logger.Info("game created", "gameID", g.ID, "status", g.Status, "private", g.Private)
	return g, nil
}

// CreateGame starts a new game in the waiting state, awaiting players to join.
func (s *GameService) CreateGame() (*Game, error) {
	return s.createGame(false)
}

func (s *GameService) CreatePrivateGame() (*Game, error) {
	return s.createGame(true)
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
	s.logger.Info("player joined", "gameID", gameID, "userID", userID, "status", g.Status)
	return g, nil
}

// JoinAnyGame finds a waiting game for the user to join, or returns an error if none are available.
func (s *GameService) JoinAnyGame(userID string) (*Game, error) {
	g, err := s.repo.FindGameToJoin(userID)
	if err != nil {
		return nil, err
	}
	return s.JoinGame(g.ID, userID)
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
	s.logger.Info("move made", "gameID", gameID, "userID", userID, "x", x, "y", y, "status", g.Status, "result", g.Result)
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
	s.logger.Info("player gave up", "gameID", gameID, "userID", userID, "winnerID", g.WinnerID)
	return g, nil
}
