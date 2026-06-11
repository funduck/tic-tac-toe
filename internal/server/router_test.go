package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/game"
)

// TestRouter verifies that routes are correctly mapped and URL parameters are extracted
func TestRouter(t *testing.T) {
	// Create a simple mock service for routing tests
	mockSvc := &mockGameService{
		createGameFunc: func() (*game.Game, error) {
			return game.NewGame("test-id"), nil
		},
		joinGameFunc: func(gameID, userID string) (*game.Game, error) {
			g := game.NewGame(gameID)
			g.UserID1 = userID
			return g, nil
		},
		makeMoveFunc: func(gameID, userID string, x, y int) (*game.Game, error) {
			g := game.NewGameInProgress(gameID, userID, "user-2")
			g.Board[x][y] = 1
			return g, nil
		},
		giveUpFunc: func(gameID, userID string) (*game.Game, error) {
			g := game.NewGameInProgress(gameID, userID, "user-2")
			g.Status = game.StatusFinished
			return g, nil
		},
		getGameFunc: func(gameID string) (*game.Game, error) {
			return game.NewGameInProgress(gameID, "user-1", "user-2"), nil
		},
	}

	handler := NewGameHandler(mockSvc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	// Create router without any additional middlewares (clean test environment)
	router := NewRouter(handler)

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "POST /api/games - create game",
			method:         http.MethodPost,
			path:           "/api/games",
			body:           CreateGameRequest{},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var g game.Game
				json.NewDecoder(resp.Body).Decode(&g)
				if g.ID != "test-id" {
					t.Errorf("expected game ID 'test-id', got '%s'", g.ID)
				}
			},
		},
		{
			name:   "POST /api/games/{gameID}/join - join game",
			method: http.MethodPost,
			path:   "/api/games/game-123/join",
			body: JoinGameRequest{
				UserID: "user-456",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var g game.Game
				json.NewDecoder(resp.Body).Decode(&g)
				if g.ID != "game-123" {
					t.Errorf("expected game ID 'game-123' (from URL), got '%s'", g.ID)
				}
				if g.UserID1 != "user-456" {
					t.Errorf("expected UserID1 'user-456', got '%s'", g.UserID1)
				}
			},
		},
		{
			name:   "POST /api/games/{gameID}/move - make move",
			method: http.MethodPost,
			path:   "/api/games/game-456/move",
			body: MoveRequest{
				UserID: "user-1",
				X:      0,
				Y:      0,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var g game.Game
				json.NewDecoder(resp.Body).Decode(&g)
				if g.ID != "game-456" {
					t.Errorf("expected game ID 'game-456' (from URL), got '%s'", g.ID)
				}
				if g.Board[0][0] != 1 {
					t.Errorf("expected board[0][0] = 1, got %d", g.Board[0][0])
				}
			},
		},
		{
			name:   "POST /api/games/{gameID}/giveup - give up",
			method: http.MethodPost,
			path:   "/api/games/game-789/giveup",
			body: GiveUpRequest{
				UserID: "user-1",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var g game.Game
				json.NewDecoder(resp.Body).Decode(&g)
				if g.ID != "game-789" {
					t.Errorf("expected game ID 'game-789' (from URL), got '%s'", g.ID)
				}
				if g.Status != game.StatusFinished {
					t.Errorf("expected status finished, got '%s'", g.Status)
				}
			},
		},
		{
			name:           "GET /api/games/{gameID} - get game",
			method:         http.MethodGet,
			path:           "/api/games/game-999",
			body:           nil,
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var g game.Game
				json.NewDecoder(resp.Body).Decode(&g)
				if g.ID != "game-999" {
					t.Errorf("expected game ID 'game-999' (from URL), got '%s'", g.ID)
				}
			},
		},
		{
			name:           "POST /api/games/{gameID}/invalid - route not found",
			method:         http.MethodPost,
			path:           "/api/games/game-123/invalid",
			body:           nil,
			expectedStatus: http.StatusNotFound,
			checkResponse:  nil,
		},
		{
			name:           "GET /api/games - wrong method",
			method:         http.MethodGet,
			path:           "/api/games",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody io.Reader
			if tt.body != nil {
				bodyBytes, _ := json.Marshal(tt.body)
				reqBody = bytes.NewBuffer(bodyBytes)
			}

			req := httptest.NewRequest(tt.method, tt.path, reqBody)
			if tt.body != nil {
				req.Header.Set("Content-Type", "application/json")
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w.Result())
			}
		})
	}
}

// TestRouterWithMiddleware verifies that custom middlewares can be injected
func TestRouterWithMiddleware(t *testing.T) {
	mockSvc := &mockGameService{
		createGameFunc: func() (*game.Game, error) {
			return game.NewGame("test-id"), nil
		},
	}

	handler := NewGameHandler(mockSvc, slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Custom middleware that adds a header
	customMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "test-value")
			next.ServeHTTP(w, r)
		})
	}

	// Create router with custom middleware
	router := NewRouter(handler, customMiddleware)

	req := httptest.NewRequest(http.MethodPost, "/api/games", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Header().Get("X-Custom-Header") != "test-value" {
		t.Errorf("expected custom header to be set by middleware")
	}
}
