package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

func TestAuthMiddleware(t *testing.T) {
	// This is a very basic test to ensure the auth middleware correctly validates tokens and user IDs.
	// More thorough testing would require mocking the TokenService and testing various scenarios (valid token, invalid token, missing token, etc.).
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")

	validToken, err := tokenService.GenerateTokens("testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	handler := AuthMiddleware(tokenService)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("valid token should allow access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		req.Header.Set("Authorization", "Bearer "+validToken.AccessToken)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", w.Code)
		}
	})

	t.Run("invalid token should deny access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d", w.Code)
		}
	})

	t.Run("invalid token format should deny access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		req.Header.Set("Authorization", "invalid-token")
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d", w.Code)
		}
	})

	t.Run("missing token should deny access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d", w.Code)
		}
	})

	t.Run("refresh token used as access token should deny access", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		req.Header.Set("Authorization", "Bearer "+validToken.RefreshToken)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401 Unauthorized, got %d", w.Code)
		}
	})

}
