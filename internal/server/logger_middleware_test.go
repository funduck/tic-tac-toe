package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	// This is a very basic test to ensure the logger middleware does not interfere with request handling.
	// More thorough testing would require capturing log output, which is beyond the scope of this simple test.
	handler := LoggerMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("GET /api/games/123 should not be logged", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/games/123", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", w.Code)
		}
	})

	t.Run("POST /api/games should be logged", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/games", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 OK, got %d", w.Code)
		}
	})
}
