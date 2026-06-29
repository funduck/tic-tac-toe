package lib

import (
	"context"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

// GameService wraps the generated API client with convenience methods
type GameService struct {
	api *openapi.GamesAPIService
}

// NewGameService creates a new GameService
func NewGameService(apiClient *openapi.APIClient) *GameService {
	return &GameService{api: apiClient.GamesAPI}
}

// CreateGame creates a new game
func (s *GameService) CreateGame(ctx context.Context, userID string, private bool) (*Game, error) {
	req := openapi.NewGameCreateGameCommand()
	req.SetPrivate(private)
	g, r, err := s.api.CreateGame(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// JoinGame joins an existing game
func (s *GameService) JoinGame(ctx context.Context, gameID, userID string) (*Game, error) {
	g, r, err := s.api.JoinGame(ctx, gameID).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// JoinAnyGame joins an existing game
func (s *GameService) JoinAnyGame(ctx context.Context, userID string) (*Game, error) {
	g, r, err := s.api.JoinAnyGame(ctx).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// MakeMove makes a move in the game
func (s *GameService) MakeMove(ctx context.Context, gameID, userID string, x, y int) (*Game, error) {
	req := openapi.NewGameMakeMoveCommand()
	req.SetX(int32(x))
	req.SetY(int32(y))
	g, r, err := s.api.MakeMove(ctx, gameID).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// GiveUp gives up the game
func (s *GameService) GiveUp(ctx context.Context, gameID, userID string) (*Game, error) {
	g, r, err := s.api.GiveUpGame(ctx, gameID).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// Quit leaves a game that is still waiting for an opponent, cancelling it.
func (s *GameService) Quit(ctx context.Context, gameID, userID string) (*Game, error) {
	g, r, err := s.api.QuitGame(ctx, gameID).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// GetGame retrieves the current game state
func (s *GameService) GetGame(ctx context.Context, gameID string) (*Game, error) {
	g, r, err := s.api.GetGame(ctx, gameID).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// GetLatestGame retrieves the authenticated user's most recent game, if any.
func (s *GameService) GetLatestGame(ctx context.Context) (*Game, error) {
	g, r, err := s.api.GetLatestGame(ctx).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}
