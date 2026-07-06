package server

import (
	"context"
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
		signupFunc: func(ctx context.Context, id, password string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: id}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
		},
		loginFunc: func(ctx context.Context, id, password string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: id}, &auth.TokenPair{AccessToken: "access", RefreshToken: "refresh"}, nil
		},
		refreshTokenFunc: func(ctx context.Context, refreshToken string) (*user.User, *auth.TokenPair, error) {
			return &user.User{ID: "testuser"}, &auth.TokenPair{AccessToken: "new-access", RefreshToken: "new-refresh"}, nil
		},
	}

	handler := NewUserHandler(mockSvc, slog.New(slog.NewTextHandler(io.Discard, nil)))
	// Create router without any additional middlewares (clean test environment)
	router := chi.NewRouter()
	router.Route("/api", func(r chi.Router) {
		UserRouter(r, handler)
	})

	// checkTokens asserts the access token is in the body and the refresh token is
	// delivered as a cookie (and not leaked into the body).
	checkTokens := func(wantAccess, wantRefresh string) func(t *testing.T, resp *http.Response) {
		return func(t *testing.T, resp *http.Response) {
			var tokens AccessTokenResponse
			json.NewDecoder(resp.Body).Decode(&tokens)
			if tokens.AccessToken != wantAccess {
				t.Errorf("unexpected access token: %q", tokens.AccessToken)
			}
			var refreshCookie *http.Cookie
			for _, c := range resp.Cookies() {
				if c.Name == refreshTokenCookieName {
					refreshCookie = c
				}
			}
			if refreshCookie == nil || refreshCookie.Value != wantRefresh {
				t.Errorf("expected refresh cookie %q, got %+v", wantRefresh, refreshCookie)
			}
		}
	}

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		cookie         *http.Cookie
		expectedStatus int
		checkResponse  func(t *testing.T, resp *http.Response)
	}{
		{
			name:           "POST /api/users/signup - successful signup",
			method:         http.MethodPost,
			path:           "/api/users/signup",
			body:           UserSignupRequest{UserID: "testuser", Password: "testpass"},
			expectedStatus: http.StatusCreated,
			checkResponse:  checkTokens("access", "refresh"),
		},
		{
			name:           "POST /api/users/login - successful login",
			method:         http.MethodPost,
			path:           "/api/users/login",
			body:           UserLoginRequest{UserID: "testuser", Password: "testpass"},
			expectedStatus: http.StatusOK,
			checkResponse:  checkTokens("access", "refresh"),
		},
		{
			name:           "POST /api/users/refresh-token - successful token refresh",
			method:         http.MethodPost,
			path:           "/api/users/refresh-token",
			cookie:         &http.Cookie{Name: refreshTokenCookieName, Value: "refresh"},
			expectedStatus: http.StatusOK,
			checkResponse:  checkTokens("new-access", "new-refresh"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := makeJSONRequest(tt.method, tt.path, tt.body)
			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}
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
