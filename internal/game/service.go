package game

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
)

// GameService processes game actions and delegates state persistence to GameRepo.
type GameService struct {
	repo             GameRepo
	logger           *slog.Logger
	joinGameMutex    sync.Mutex // protects JoinGame
	joinAnyGameMutex sync.Mutex // protects JoinAnyGame
}

func NewGameService(repo GameRepo, logger *slog.Logger) *GameService {
	return &GameService{repo: repo, logger: logger}
}

func (s *GameService) createGame(ctx context.Context, private bool) (*Game, error) {
	id, err := uuid.NewV7() // V7 for lexicographically sortable IDs
	if err != nil {
		return nil, err
	}
	g := NewGame(id.String())
	if private {
		g.Private = true
	}
	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("game created", "gameID", g.ID, "status", g.Status, "private", g.Private)
	return g, nil
}

// CreateGame starts a new game in the waiting state, awaiting players to join.
func (s *GameService) CreateGame(ctx context.Context) (*Game, error) {
	return s.createGame(ctx, false)
}

func (s *GameService) CreatePrivateGame(ctx context.Context) (*Game, error) {
	return s.createGame(ctx, true)
}

// JoinGame adds a player to an existing waiting game. The first player to join becomes UserID1, the second becomes UserID2.
// If both players have joined, the game status changes to in_progress and UserID1 moves first.
func (s *GameService) JoinGame(ctx context.Context, gameID, userID string) (*Game, error) {
	// To make this method concurrency-safe we need to wrap the code into a "transaction"
	// With SQL database we'd use a transaction with "SELECT ... FOR UPDATE"
	s.joinGameMutex.Lock()
	defer s.joinGameMutex.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	g, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if err := g.Join(userID); err != nil {
		return nil, err
	}
	g.Touch(userID, time.Now())
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}

	s.logger.Info("player joined", "gameID", gameID, "userID", userID, "status", g.Status)
	return g, nil
}

// JoinAnyGame finds a waiting game for the user to join, or returns an error if none are available.
func (s *GameService) JoinAnyGame(ctx context.Context, userID string) (*Game, error) {
	// To make this method concurrency-safe we need to wrap the code into a "transaction"
	// With SQL database we'd use a transaction with "SELECT ... FOR UPDATE"
	// We need another mutex here to avoid picking the same game for multiple users concurrently
	s.joinAnyGameMutex.Lock()
	defer s.joinAnyGameMutex.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	g, err := s.repo.FindGameToJoin(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.JoinGame(ctx, g.ID, userID)
}

// GetLatestGameForUser retrieves the most recent game for a given user.
func (s *GameService) GetLatestGameForUser(ctx context.Context, userID string) (*Game, error) {
	return s.repo.FindLatestForUser(ctx, userID)
}

// GetGame retrieves a game by ID and records the requesting user's presence.
func (s *GameService) GetGame(ctx context.Context, gameID, userID string) (*Game, error) {
	g, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if g.Touch(userID, time.Now()) {
		if err := s.repo.Update(ctx, g); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// MakeMove applies a move on behalf of userID and persists the updated state.
func (s *GameService) MakeMove(ctx context.Context, gameID, userID string, x, y int) (*Game, error) {
	g, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if err := g.MakeMove(userID, x, y); err != nil {
		return nil, err
	}
	g.Touch(userID, time.Now())
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("move made", "gameID", gameID, "userID", userID, "x", x, "y", y, "status", g.Status, "result", g.Result)
	return g, nil
}

// GiveUp concedes the game for userID, awarding the win to the other player.
func (s *GameService) GiveUp(ctx context.Context, gameID, userID string) (*Game, error) {
	g, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if err := g.GiveUp(userID); err != nil {
		return nil, err
	}
	g.Touch(userID, time.Now())
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("player gave up", "gameID", gameID, "userID", userID, "winnerID", g.WinnerID)
	return g, nil
}

// Quit cancels a game that is still waiting for an opponent, on behalf of userID.
func (s *GameService) Quit(ctx context.Context, gameID, userID string) (*Game, error) {
	g, err := s.repo.FindByID(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if err := g.Quit(userID); err != nil {
		return nil, err
	}
	g.Touch(userID, time.Now())
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("player quit", "gameID", gameID, "userID", userID, "status", g.Status)
	return g, nil
}
