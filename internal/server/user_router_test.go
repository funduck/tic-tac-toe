package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/user"
	"github.com/go-chi/chi/v5"
)

func TestUserRouter(t *testing.T) {
	// Create a simple mock service for routing tests
	mockSvc := &mockUserService{
		signupFunc: func(id, password string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: id}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
		},
		loginFunc: func(id, password string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: id}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
		},
		refreshTokenFunc: func(id, refreshToken string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: id}, &auth.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
		},
	}

	handler := NewUserHandler(mockSvc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	// Create router without any additional middlewares (clean test environment)
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		UserRouter(r, handler)
	})

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "POST /api/users/signup - successful signup",
			method:         http.MethodPost,
			path:           "/api/users/signup",
			body:           UserSignupRequest{UserID: "testuser", Password: "testpass"},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var tokens auth.TokenPair
				json.NewDecoder(resp.Body).Decode(&tokens)
				if tokens.AccessToken != "access" || tokens.RefreshToken != "refresh" {
					t.Errorf("unexpected tokens: %+v", tokens)
				}
			},
		},
		{
			name:           "POST /api/users/login - successful login",
			method:         http.MethodPost,
			path:           "/api/users/login",
			body:           UserLoginRequest{UserID: "testuser", Password: "testpass"},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var tokens auth.TokenPair
				json.NewDecoder(resp.Body).Decode(&tokens)
				if tokens.AccessToken != "access" || tokens.RefreshToken != "refresh" {
					t.Errorf("unexpected tokens: %+v", tokens)
				}
			},
		},
		{
			name:           "POST /api/users/refresh-token - successful token refresh",
			method:         http.MethodPost,
			path:           "/api/users/refresh-token",
			body:           UserRefreshTokenRequest{UserID: "testuser", RefreshToken: "refresh"},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *http.Response) {
				var tokens auth.TokenPair
				json.NewDecoder(resp.Body).Decode(&tokens)
				if tokens.AccessToken != "new-access" || tokens.RefreshToken != "new-refresh" {
					t.Errorf("unexpected tokens: %+v", tokens)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := makeJSONRequest(tt.method, tt.path, tt.body)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.Code)
			}
			if tt.checkResponse != nil {
				tt.checkResponse(t, resp.Result())
			}
		})
	}
}
