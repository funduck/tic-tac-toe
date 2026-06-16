package lib

import (
	"context"
	"time"

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
func (gs *GameService) CreateGame(ctx context.Context, userID string, private bool) (*Game, error) {
	req := openapi.NewServerCreateGameRequest()
	req.SetUserID(userID)
	req.SetPrivate(private)
	g, r, err := gs.api.CreateGame(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// JoinGame joins an existing game
func (gs *GameService) JoinGame(ctx context.Context, gameID, userID string) (*Game, error) {
	req := openapi.NewServerJoinGameRequest()
	req.SetUserID(userID)
	g, r, err := gs.api.JoinGame(ctx, gameID).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// JoinAnyGame joins an existing game
func (gs *GameService) JoinAnyGame(ctx context.Context, userID string) (*Game, error) {
	req := openapi.NewServerJoinAnyGameRequest()
	req.SetUserID(userID)
	g, r, err := gs.api.JoinAnyGame(ctx).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// MakeMove makes a move in the game
func (gs *GameService) MakeMove(ctx context.Context, gameID, userID string, x, y int) (*Game, error) {
	req := openapi.NewServerMoveRequest()
	req.SetUserID(userID)
	req.SetX(int32(x))
	req.SetY(int32(y))
	g, r, err := gs.api.MakeMove(ctx, gameID).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// GiveUp gives up the game
func (gs *GameService) GiveUp(ctx context.Context, gameID, userID string) (*Game, error) {
	req := openapi.NewServerGiveUpRequest()
	req.SetUserID(userID)
	g, r, err := gs.api.GiveUpGame(ctx, gameID).Request(*req).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// GetGame retrieves the current game state
func (gs *GameService) GetGame(ctx context.Context, gameID string) (*Game, error) {
	g, r, err := gs.api.GetGame(ctx, gameID).Execute()
	if err != nil {
		return nil, parseError(r, err)
	}
	return WrapGame(g), nil
}

// PollUntil polls the game until the predicate returns true or an error occurs
// Retries up to maxRetries times with exponential backoff on errors
func (gs *GameService) PollUntil(ctx context.Context, gameID string, predicate func(*Game) bool, maxRetries int) (*Game, error) {
	retries := 0
	for {
		g, err := gs.GetGame(ctx, gameID)
		if err != nil {
			retries++
			if retries >= maxRetries {
				return nil, err
			}
			// Exponential backoff
			time.Sleep(time.Duration(retries) * time.Second)
			continue
		}
		// Reset retries on success
		retries = 0

		if predicate(g) {
			return g, nil
		}
		time.Sleep(time.Second)
	}
}
