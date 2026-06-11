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
	createGameFunc        func() (*game.Game, error)
	createPrivateGameFunc func() (*game.Game, error)
	joinGameFunc          func(gameID, userID string) (*game.Game, error)
	joinAnyGameFunc       func(userID string) (*game.Game, error)
	makeMoveFunc          func(gameID, userID string, x, y int) (*game.Game, error)
	giveUpFunc            func(gameID, userID string) (*game.Game, error)
	getGameFunc           func(gameID string) (*game.Game, error)
}

func (m *mockGameService) CreateGame() (*game.Game, error) {
	if m.createGameFunc != nil {
		return m.createGameFunc()
	}
	return nil, nil
}

func (m *mockGameService) CreatePrivateGame() (*game.Game, error) {
	if m.createPrivateGameFunc != nil {
		return m.createPrivateGameFunc()
	}
	return nil, nil
}

func (m *mockGameService) JoinGame(gameID, userID string) (*game.Game, error) {
	if m.joinGameFunc != nil {
		return m.joinGameFunc(gameID, userID)
	}
	return nil, nil
}

func (m *mockGameService) JoinAnyGame(userID string) (*game.Game, error) {
	if m.joinAnyGameFunc != nil {
		return m.joinAnyGameFunc(userID)
	}
	return nil, nil
}

func (m *mockGameService) MakeMove(gameID, userID string, x, y int) (*game.Game, error) {
	if m.makeMoveFunc != nil {
		return m.makeMoveFunc(gameID, userID, x, y)
	}
	return nil, nil
}

func (m *mockGameService) GiveUp(gameID, userID string) (*game.Game, error) {
	if m.giveUpFunc != nil {
		return m.giveUpFunc(gameID, userID)
	}
	return nil, nil
}

func (m *mockGameService) GetGame(gameID string) (*game.Game, error) {
	if m.getGameFunc != nil {
		return m.getGameFunc(gameID)
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

// executeRequest executes a request and returns the response recorder
func executeRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler(w, req)
	return w
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
		requestBody    CreateGameRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "create game without user",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func() (*game.Game, error) {
						return game.NewGame("test-game-id-123"), nil
					},
				}
			},
			requestBody:    CreateGameRequest{},
			expectedStatus: http.StatusCreated,
			checkResponse: checkGameMatches(game.Game{
				ID:     "test-game-id-123",
				Status: game.StatusWaiting,
			}),
		},
		{
			name: "create game with user",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func() (*game.Game, error) {
						return game.NewGame("test-game-id-123"), nil
					},
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						g := game.NewGame(gameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			requestBody: CreateGameRequest{
				UserID: "user-123",
			},
			expectedStatus: http.StatusCreated,
			checkResponse: checkGameMatches(game.Game{
				ID:      "test-game-id-123",
				Status:  game.StatusWaiting,
				UserID1: "user-123",
			}),
		},
		{
			name: "create private game with user",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createPrivateGameFunc: func() (*game.Game, error) {
						g := game.NewGame("test-game-id-123")
						g.Private = true
						return g, nil
					},
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						g := game.NewGame(gameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			requestBody: CreateGameRequest{
				UserID:  "user-123",
				Private: true,
			},
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
					createGameFunc: func() (*game.Game, error) {
						return nil, errors.New("database error")
					},
				}
			},
			requestBody:    CreateGameRequest{},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  checkErrorResponse("database error"),
		},
		{
			name: "join game service error",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					createGameFunc: func() (*game.Game, error) {
						return game.NewGame("test-game-id-123"), nil
					},
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						return nil, errors.New("join error")
					},
				}
			},
			requestBody: CreateGameRequest{
				UserID: "user-123",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse:  checkErrorResponse("join error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
			)

			req := makeJSONRequest(http.MethodPost, "/api/games", tt.requestBody)
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
		requestBody    JoinGameRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully join game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						g := game.NewGame(gameID)
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			gameID: "game-123",
			requestBody: JoinGameRequest{
				UserID: "user-456",
			},
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
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			gameID: "nonexistent",
			requestBody: JoinGameRequest{
				UserID: "user-456",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
		{
			name: "game not waiting",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinGameFunc: func(gameID, userID string) (*game.Game, error) {
						return nil, game.ErrGameNotWaiting
					},
				}
			},
			gameID: "game-123",
			requestBody: JoinGameRequest{
				UserID: "user-456",
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrGameNotWaiting.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
			)

			req := makeJSONRequest(http.MethodPost, "/api/games/"+tt.gameID+"/join", tt.requestBody)
			// Add URL param to context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

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
		requestBody    JoinAnyGameRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully join any game",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					joinAnyGameFunc: func(userID string) (*game.Game, error) {
						g := game.NewGame("game-123")
						g.UserID1 = userID
						return g, nil
					},
				}
			},
			requestBody: JoinAnyGameRequest{
				UserID: "user-456",
			},
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
					joinAnyGameFunc: func(userID string) (*game.Game, error) {
						return nil, game.ErrGameNotFound
					},
				}
			},
			requestBody: JoinAnyGameRequest{
				UserID: "user-456",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse:  checkErrorResponse(game.ErrGameNotFound.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
			)

			req := makeJSONRequest(http.MethodPost, "/api/games/join", tt.requestBody)
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
		requestBody    MoveRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully make move",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					makeMoveFunc: func(gameID, userID string, x, y int) (*game.Game, error) {
						g := game.NewGameInProgress(gameID, userID, "user-2")
						g.Board[x][y] = 1
						g.CurrentPlayerID = "user-2"
						return g, nil
					},
				}
			},
			gameID: "game-123",
			requestBody: MoveRequest{
				UserID: "user-1",
				X:      0,
				Y:      0,
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
					makeMoveFunc: func(gameID, userID string, x, y int) (*game.Game, error) {
						return nil, game.ErrNotYourTurn
					},
				}
			},
			gameID: "game-123",
			requestBody: MoveRequest{
				UserID: "user-2",
				X:      0,
				Y:      0,
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrNotYourTurn.Error()),
		},
		{
			name: "cell occupied",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					makeMoveFunc: func(gameID, userID string, x, y int) (*game.Game, error) {
						return nil, game.ErrCellOccupied
					},
				}
			},
			gameID: "game-123",
			requestBody: MoveRequest{
				UserID: "user-1",
				X:      1,
				Y:      1,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse:  checkErrorResponse(game.ErrCellOccupied.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
			)

			req := makeJSONRequest(http.MethodPost, "/api/games/"+tt.gameID+"/move", tt.requestBody)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

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
		requestBody    GiveUpRequest
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "successfully give up",
			mockSetup: func() *mockGameService {
				return &mockGameService{
					giveUpFunc: func(gameID, userID string) (*game.Game, error) {
						g := game.NewGameInProgress(gameID, userID, "user-2")
						g.Status = game.StatusFinished
						g.Result = game.ResultWin
						g.WinnerID = "user-2"
						return g, nil
					},
				}
			},
			gameID: "game-123",
			requestBody: GiveUpRequest{
				UserID: "user-1",
			},
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
					giveUpFunc: func(gameID, userID string) (*game.Game, error) {
						return nil, game.ErrNotInGame
					},
				}
			},
			gameID: "game-123",
			requestBody: GiveUpRequest{
				UserID: "user-999",
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  checkErrorResponse(game.ErrNotInGame.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewGameHandler(
				tt.mockSetup(),
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
			)

			req := makeJSONRequest(http.MethodPost, "/api/games/"+tt.gameID+"/giveup", tt.requestBody)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("gameID", tt.gameID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := executeRequest(h.GiveUp, req)

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
					getGameFunc: func(gameID string) (*game.Game, error) {
						g := game.NewGameInProgress(gameID, "user-1", "user-2")
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
					getGameFunc: func(gameID string) (*game.Game, error) {
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
				slog.New(slog.NewTextHandler(io.Discard, nil)), // silent logger for tests
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
