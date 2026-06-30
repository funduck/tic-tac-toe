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
	afkTimeout       time.Duration
}

func NewGameService(repo GameRepo, logger *slog.Logger, afkTimeout time.Duration) *GameService {
	return &GameService{repo: repo, logger: logger, afkTimeout: afkTimeout}
}

type CreateGameCommand struct {
	Private bool
}

type JoinGameCommand struct {
	GameID string
}

type GetGameCommand struct {
	GameID string
}

type MakeMoveCommand struct {
	GameID string
	X      int
	Y      int
}

type QuitCommand struct {
	GameID string
}

type GiveUpCommand struct {
	GameID string
}

// CreateGame starts a new game in the waiting state, awaiting players to join.
func (s *GameService) CreateGame(ctx context.Context, userId string, cmd CreateGameCommand) (*Game, error) {
	id, err := uuid.NewV7() // V7 for lexicographically sortable IDs
	if err != nil {
		return nil, err
	}
	g := NewGame(id.String())
	if cmd.Private {
		g.Private = true
	}
	if err := s.repo.Create(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("game created", "userId", userId, "gameID", g.ID, "status", g.Status, "private", g.Private)
	return g, nil
}

// JoinGame adds a player to an existing waiting game. The first player to join becomes UserID1, the second becomes UserID2.
// If both players have joined, the game status changes to in_progress and UserID1 moves first.
func (s *GameService) JoinGame(ctx context.Context, userID string, cmd JoinGameCommand) (*Game, error) {
	// Simple mutex to avoid two users joining the same game concurrently
	s.joinGameMutex.Lock()
	defer s.joinGameMutex.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	g, err := s.repo.FindByID(ctx, cmd.GameID)
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

	s.logger.Info("player joined", "userID", userID, "gameID", cmd.GameID, "status", g.Status)
	return g, nil
}

// JoinAnyGame finds a waiting game for the user to join, or returns an error if none are available.
func (s *GameService) JoinAnyGame(ctx context.Context, userID string) (*Game, error) {
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

	return s.JoinGame(ctx, userID, JoinGameCommand{GameID: g.ID})
}

// GetLatestGameForUser retrieves the most recent game for a given user.
func (s *GameService) GetLatestGameForUser(ctx context.Context, userID string) (*Game, error) {
	return s.repo.FindLatestForUser(ctx, userID)
}

// GetGame retrieves a game by ID and records the requesting user's presence.
func (s *GameService) GetGame(ctx context.Context, userID string, cmd GetGameCommand) (*Game, error) {
	g, err := s.repo.FindByID(ctx, cmd.GameID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	changed := g.Touch(userID, now)
	if g.ForfeitIfOpponentAFK(userID, now, s.afkTimeout) {
		changed = true
		s.logger.Info("opponent AFK, auto-win", "gameID", cmd.GameID, "winnerID", g.WinnerID)
	}
	if changed {
		if err := s.repo.Update(ctx, g); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// MakeMove applies a move on behalf of userID and persists the updated state.
func (s *GameService) MakeMove(ctx context.Context, userID string, cmd MakeMoveCommand) (*Game, error) {
	g, err := s.repo.FindByID(ctx, cmd.GameID)
	if err != nil {
		return nil, err
	}
	if err := g.MakeMove(userID, cmd.X, cmd.Y); err != nil {
		return nil, err
	}
	g.Touch(userID, time.Now())
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	s.logger.Info("move made", "userID", userID, "gameID", cmd.GameID, "x", cmd.X, "y", cmd.Y, "status", g.Status, "result", g.Result)
	return g, nil
}

// GiveUp concedes the game for userID, awarding the win to the other player.
func (s *GameService) GiveUp(ctx context.Context, userID string, cmd GiveUpCommand) (*Game, error) {
	g, err := s.repo.FindByID(ctx, cmd.GameID)
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
	s.logger.Info("player gave up", "userID", userID, "gameID", cmd.GameID, "winnerID", g.WinnerID)
	return g, nil
}

// Quit cancels a game that is still waiting for an opponent, on behalf of userID.
func (s *GameService) Quit(ctx context.Context, userID string, cmd QuitCommand) (*Game, error) {
	g, err := s.repo.FindByID(ctx, cmd.GameID)
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
	s.logger.Info("player quit", "userID", userID, "gameID", cmd.GameID, "status", g.Status)
	return g, nil
}
