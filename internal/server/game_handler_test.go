package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/game"
	"github.com/go-chi/chi/v5"
)

type mockGameService struct {
	createGameFunc    func(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error)
	joinGameFunc      func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error)
	joinAnyGameFunc   func(ctx context.Context, userID string) (*game.Game, error)
	makeMoveFunc      func(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error)
	giveUpFunc        func(ctx context.Context, userID string, cmd game.GiveUpCommand) (*game.Game, error)
	quitFunc          func(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error)
	getGameFunc       func(ctx context.Context, userID string, cmd game.GetGameCommand) (*game.Game, error)
	getLatestGameFunc func(ctx context.Context, userID string) (*game.Game, error)
}

func (m *mockGameService) CreateGame(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error) {
	if m.createGameFunc != nil {
		return m.createGameFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) JoinGame(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
	if m.joinGameFunc != nil {
		return m.joinGameFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) JoinAnyGame(ctx context.Context, userID string) (*game.Game, error) {
	if m.joinAnyGameFunc != nil {
		return m.joinAnyGameFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockGameService) MakeMove(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error) {
	if m.makeMoveFunc != nil {
		return m.makeMoveFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) GiveUp(ctx context.Context, userID string, cmd game.GiveUpCommand) (*game.Game, error) {
	if m.giveUpFunc != nil {
		return m.giveUpFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) Quit(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error) {
	if m.quitFunc != nil {
		return m.quitFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) GetGame(ctx context.Context, userID string, cmd game.GetGameCommand) (*game.Game, error) {
	if m.getGameFunc != nil {
		return m.getGameFunc(ctx, userID, cmd)
	}
	return nil, nil
}

func (m *mockGameService) GetLatestGameForUser(ctx context.Context, userID string) (*game.Game, error) {
	if m.getLatestGameFunc != nil {
		return m.getLatestGameFunc(ctx, userID)
	}
	return nil, nil
}

// Test utility functions

// makeJSONRequest creates an HTTP request with JSON body
func makeJSONRequest(method, url string, body interface{}) *http.Request {
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// withUserID injects a userID into the request context (simulates AuthMiddleware)
func withUserID(req *http.Request, userID string) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), userIDContextKey, userID))
}

// executeRequest executes a request and returns the response recorder
func executeRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

// findCookie returns the cookie with the given name from a recorded response, or nil.
func findCookie(w *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// checkStatusCode returns a function that checks the HTTP status code
func checkStatusCode(expectedStatus int) func(t *testing.T, w *httptest.ResponseRecorder) {
	return func(t *testing.T, w *httptest.ResponseRecorder) {
		if w.Code != expectedStatus {
			t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
		}
	}
}

// checkGameMatches returns a function that validates the response game matches expected values
func checkGameMatches(expected game.Game) func(t *testing.T, body []byte) {
	return func(t *testing.T, body []byte) {
		var g game.Game
		if err := json.Unmarshal(body, &g); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
		if expected.ID != "" && g.ID != expected.ID {
			t.Errorf("expected game ID '%s', got '%s'", expected.ID, g.ID)
		}
		if expected.Status != "" && g.Status != expected.Status {
			t.Errorf("expected status '%s', got '%s'", expected.Status, g.Status)
		}
		if expected.UserID1 != "" && g.UserID1 != expected.UserID1 {
			t.Errorf("expected UserID1 '%s', got '%s'", expected.UserID1, g.UserID1)
		}
		if expected.UserID2 != "" && g.UserID2 != expected.UserID2 {
			t.Errorf("expected UserID2 '%s', got '%s'", expected.UserID2, g.UserID2)
		}
		if expected.CurrentPlayerID != "" && g.CurrentPlayerID != expected.CurrentPlayerID {
			t.Errorf("expected CurrentPlayerID '%s', got '%s'", expected.CurrentPlayerID, g.CurrentPlayerID)
		}
		if expected.Result != "" && g.Result != expected.Result {
			t.Errorf("expected Result '%s', got '%s'", expected.Result, g.Result)
		}
		if expected.WinnerID != "" && g.WinnerID != expected.WinnerID {
			t.Errorf("expected WinnerID '%s', got '%s'", expected.WinnerID, g.WinnerID)
		}
	}
}

// checkErrorResponse validates the error response contains expected message
func checkErrorResponse(expectedError string) func(t *testing.T, body []byte) {
	return func(t *testing.T, body []byte) {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			t.Fatalf("failed to unmarshal error response: %v", err)
		}
		if errResp.Error != expectedError {
			t.Errorf("expected error '%s', got '%s'", expectedError, errResp.Error)
		}
	}
}

func TestGameHandler_CreateGame(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		userID         string
		requestBody    game.CreateGameCommand
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "create game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error) {
						return game.NewGame("test-game-id-123"), nil
					},
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						g := game.NewGame(cmd.GameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			userID:         "user-123",
			requestBody:    game.CreateGameCommand{},
			expectedStatus: http.StatusCreated,
			checkResponse: checkGameMatches(game.Game{
				ID:      "test-game-id-123",
				Status:  game.StatusWaiting,
				UserID1: "user-123",
			}),
		},
		{
			name: "create private game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error) {
						g := game.NewGame("test-game-id-123")
						g.Private = cmd.Private
						return g, nil
					},
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						g := game.NewGame(cmd.GameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			userID:         "user-123",
			requestBody:    game.CreateGameCommand{Private: true},
			expectedStatus: http.StatusCreated,
			checkResponse: checkGameMatches(game.Game{
				ID:      "test-game-id-123",
				Status:  game.StatusWaiting,
				UserID1: "user-123",
				Private: true,
			}),
		},
		{
			name: "create game service error",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error) {
						return nil, errors.New("database error")
					},
				}
			},
			userID:         "user-123",
			requestBody:    game.CreateGameCommand{},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  checkErrorResponse("database error"),
		},
		{
			name: "join game service error",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func(ctx context.Context, userID string, cmd game.CreateGameCommand) (*game.Game, error) {
						return game.NewGame("test-game-id-123"), nil
					},
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						return nil, errors.New("join error")
					},
				}
			},
			userID:         "user-123",
			requestBody:    game.CreateGameCommand{},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  checkErrorResponse("join error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := makeJSONRequest(http.MethodPost, "/api/games", tt.requestBody)
			req = withUserID(req, tt.userID)
			w := executeRequest(h.CreateGame, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_JoinGame(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		gameID         string
		userID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully join game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						g := game.NewGame(cmd.GameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-456",
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:      "game-123",
				Status:  game.StatusWaiting,
				UserID1: "user-456",
			}),
		},
		{
			name: "join game not found",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			gameID:         "nonexistent",
			userID:         "user-456",
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
		{
			name: "game not waiting",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinGameFunc: func(ctx context.Context, userID string, cmd game.JoinGameCommand) (*game.Game, error) {
						return nil, game.ErrGameNotWaiting
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-456",
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrGameNotWaiting.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/games/"+tt.gameID+"/join", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = withUserID(req, tt.userID)

			w := executeRequest(h.JoinGame, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_JoinAnyGame(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		userID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully join any game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinAnyGameFunc: func(ctx context.Context, userID string) (*game.Game, error) {
						g := game.NewGame("game-123")
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			userID:         "user-456",
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:      "game-123",
				Status:  game.StatusWaiting,
				UserID1: "user-456",
			}),
		},
		{
			name: "no game found",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinAnyGameFunc: func(ctx context.Context, userID string) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			userID:         "user-456",
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/games/join", nil)
			req = withUserID(req, tt.userID)
			w := executeRequest(h.JoinAnyGame, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_MakeMove(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		gameID         string
		userID         string
		requestBody    game.MakeMoveCommand
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully make move",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					makeMoveFunc: func(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error) {
						g := game.NewGameInProgress(cmd.GameID, userID, "user-2")
						g.Board[cmd.X][cmd.Y] = 1
						g.CurrentPlayerID = "user-2"
						return g, nil
					},
				}
			},
			gameID: "game-123",
			userID: "user-1",
			requestBody: game.MakeMoveCommand{
				X: 0,
				Y: 0,
			},
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:              "game-123",
				Status:          game.StatusInProgress,
				UserID1:         "user-1",
				UserID2:         "user-2",
				CurrentPlayerID: "user-2",
			}),
		},
		{
			name: "not your turn",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					makeMoveFunc: func(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error) {
						return nil, game.ErrNotYourTurn
					},
				}
			},
			gameID: "game-123",
			userID: "user-2",
			requestBody: game.MakeMoveCommand{
				X: 0,
				Y: 0,
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrNotYourTurn.Error()),
		},
		{
			name: "cell occupied",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					makeMoveFunc: func(ctx context.Context, userID string, cmd game.MakeMoveCommand) (*game.Game, error) {
						return nil, game.ErrCellOccupied
					},
				}
			},
			gameID: "game-123",
			userID: "user-1",
			requestBody: game.MakeMoveCommand{
				X: 1,
				Y: 1,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  checkErrorResponse(game.ErrCellOccupied.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := makeJSONRequest(http.MethodPost, "/api/games/"+tt.gameID+"/move", tt.requestBody)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = withUserID(req, tt.userID)

			w := executeRequest(h.MakeMove, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_GiveUp(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		gameID         string
		userID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully give up",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					giveUpFunc: func(ctx context.Context, userID string, cmd game.GiveUpCommand) (*game.Game, error) {
						g := game.NewGameInProgress(cmd.GameID, userID, "user-2")
						g.Status = game.StatusFinished
						g.Result = game.ResultWin
						g.WinnerID = "user-2"
						return g, nil
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-1",
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:       "game-123",
				Status:   game.StatusFinished,
				Result:   game.ResultWin,
				WinnerID: "user-2",
			}),
		},
		{
			name: "player not in game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					giveUpFunc: func(ctx context.Context, userID string, cmd game.GiveUpCommand) (*game.Game, error) {
						return nil, game.ErrNotInGame
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-999",
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrNotInGame.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/games/"+tt.gameID+"/giveup", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = withUserID(req, tt.userID)

			w := executeRequest(h.GiveUp, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_Quit(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		gameID         string
		userID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully quit waiting game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					quitFunc: func(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error) {
						g := game.NewGame(cmd.GameID)
						g.UserID1 = userID
						g.Status = game.StatusCancelled
						return g, nil
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-1",
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:      "game-123",
				Status:  game.StatusCancelled,
				UserID1: "user-1",
			}),
		},
		{
			name: "game not waiting",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					quitFunc: func(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error) {
						return nil, game.ErrGameNotWaiting
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-1",
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrGameNotWaiting.Error()),
		},
		{
			name: "player not in game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					quitFunc: func(ctx context.Context, userID string, cmd game.QuitCommand) (*game.Game, error) {
						return nil, game.ErrNotInGame
					},
				}
			},
			gameID:         "game-123",
			userID:         "user-999",
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrNotInGame.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodPost, "/api/games/"+tt.gameID+"/quit", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			req = withUserID(req, tt.userID)

			w := executeRequest(h.Quit, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_GetGame(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		gameID         string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully get game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					getGameFunc: func(ctx context.Context, userID string, cmd game.GetGameCommand) (*game.Game, error) {
						g := game.NewGameInProgress(cmd.GameID, "user-1", "user-2")
						return g, nil
					},
				}
			},
			gameID:         "game-123",
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:              "game-123",
				Status:          game.StatusInProgress,
				UserID1:         "user-1",
				UserID2:         "user-2",
				CurrentPlayerID: "user-1",
			}),
		},
		{
			name: "game not found",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					getGameFunc: func(ctx context.Context, userID string, cmd game.GetGameCommand) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			gameID:         "nonexistent",
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodGet, "/api/games/"+tt.gameID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := executeRequest(h.GetGame, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}

func TestGameHandler_GetLatestGame(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func() *mockGameService
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully get latest game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					getLatestGameFunc: func(ctx context.Context, userID string) (*game.Game, error) {
						return game.NewGameInProgress("game-123", userID, "user-2"), nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: checkGameMatches(game.Game{
				ID:              "game-123",
				Status:          game.StatusInProgress,
				UserID1:         "test-user",
				UserID2:         "user-2",
				CurrentPlayerID: "test-user",
			}),
		},
		{
			name: "no game for user",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					getLatestGameFunc: func(ctx context.Context, userID string) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)),
			)

			req := httptest.NewRequest(http.MethodGet, "/api/games", nil)
			req = req.WithContext(context.WithValue(req.Context(), userIDContextKey, "test-user"))

			w := executeRequest(h.GetLatestGame, req)

			checkStatusCode(tt.expectedStatus)(t, w)
			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Body.Bytes())
			}
		})
	}
}
