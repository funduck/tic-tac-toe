package game

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// GameService processes game actions and delegates state persistence to GameRepo.
type GameService struct {
	repo       GameRepo
	logger     *slog.Logger
	afkTimeout time.Duration
}

func NewGameService(repo GameRepo, logger *slog.Logger) *GameService {
	return &GameService{repo: repo, logger: logger, afkTimeout: 10 * time.Second}
}

func (s *GameService) SetAfkTimeout(duration time.Duration) {
	s.afkTimeout = duration
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

var (
	// errSkipUpdate is returned by an updateGame mutate func to report that the
	// game did not change and does not need to be persisted.
	errSkipUpdate = errors.New("no changes to persist")

	maxRetries = 3

	// errMaxRetries is returned by updateGame when the maximum number of retries is reached.
	errMaxRetries = errors.New("max retries reached")
)

// updateGame runs a read-modify-write cycle on a game with optimistic locking.
func (s *GameService) updateGame(ctx context.Context, gameID string, mutate func(g *Game) error) (*Game, error) {
	retries := 0
	for retries < maxRetries {
		g, err := s.repo.FindByID(ctx, gameID)
		if err != nil {
			return nil, err
		}
		if err := mutate(g); err != nil {
			if errors.Is(err, errSkipUpdate) {
				return g, nil
			}
			return nil, err
		}
		err = s.repo.Update(ctx, g)
		if err == nil {
			return g, nil
		}
		// return on any error other than version conflict, which indicates a concurrent update
		if !errors.Is(err, ErrVersionConflict) {
			return nil, err
		}
		// return on context error
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		retries++
	}
	return nil, errMaxRetries
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
	g, err := s.updateGame(ctx, cmd.GameID, func(g *Game) error {
		if err := g.Join(userID); err != nil {
			return err
		}
		g.Touch(userID, time.Now())
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info("player joined", "userID", userID, "gameID", cmd.GameID, "status", g.Status)
	return g, nil
}

// JoinAnyGame finds a waiting game for the user to join, or returns an error if none are available.
func (s *GameService) JoinAnyGame(ctx context.Context, userID string) (*Game, error) {
	for {
		g, err := s.repo.FindGameToJoin(ctx, userID)
		if err != nil {
			return nil, err
		}
		g, err = s.JoinGame(ctx, userID, JoinGameCommand{GameID: g.ID})
		if err == nil {
			return g, nil
		}
		// A concurrent player filled this game first; look for another one.
		if !errors.Is(err, ErrGameNotWaiting) {
			return nil, err
		}
		// return on context error
		if err := ctx.Err(); err != nil {
			return nil, err
		}
	}
}

// GetLatestGameForUser retrieves the most recent game for a given user.
func (s *GameService) GetLatestGameForUser(ctx context.Context, userID string) (*Game, error) {
	return s.repo.FindLatestForUser(ctx, userID)
}

// GetGame retrieves a game by ID and records the requesting user's presence.
func (s *GameService) GetGame(ctx context.Context, userID string, cmd GetGameCommand) (*Game, error) {
	return s.updateGame(ctx, cmd.GameID, func(g *Game) error {
		now := time.Now()
		changed := g.Touch(userID, now)
		if g.ForfeitIfOpponentAFK(userID, now, s.afkTimeout) {
			changed = true
			s.logger.Info("opponent AFK, auto-win", "gameID", cmd.GameID, "winnerID", g.WinnerID)
		}
		if !changed {
			return errSkipUpdate
		}
		return nil
	})
}

// MakeMove applies a move on behalf of userID and persists the updated state.
func (s *GameService) MakeMove(ctx context.Context, userID string, cmd MakeMoveCommand) (*Game, error) {
	g, err := s.updateGame(ctx, cmd.GameID, func(g *Game) error {
		if err := g.MakeMove(userID, cmd.X, cmd.Y); err != nil {
			return err
		}
		g.Touch(userID, time.Now())
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info("move made", "userID", userID, "gameID", cmd.GameID, "x", cmd.X, "y", cmd.Y, "status", g.Status, "result", g.Result)
	return g, nil
}

// GiveUp concedes the game for userID, awarding the win to the other player.
func (s *GameService) GiveUp(ctx context.Context, userID string, cmd GiveUpCommand) (*Game, error) {
	g, err := s.updateGame(ctx, cmd.GameID, func(g *Game) error {
		if err := g.GiveUp(userID); err != nil {
			return err
		}
		g.Touch(userID, time.Now())
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info("player gave up", "userID", userID, "gameID", cmd.GameID, "winnerID", g.WinnerID)
	return g, nil
}

// Quit cancels a game that is still waiting for an opponent, on behalf of userID.
func (s *GameService) Quit(ctx context.Context, userID string, cmd QuitCommand) (*Game, error) {
	g, err := s.updateGame(ctx, cmd.GameID, func(g *Game) error {
		if err := g.Quit(userID); err != nil {
			return err
		}
		g.Touch(userID, time.Now())
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info("player quit", "userID", userID, "gameID", cmd.GameID, "status", g.Status)
	return g, nil
}
