package server

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5/middleware"
)

var regexGetGame = regexp.MustCompile(`^/api/games/[^/]+$`)

// loggerMiddleware is a middleware that logs all requests except GET /games/{gameID} (polling)
func loggerMiddleware(next http.Handler) http.Handler {
	logger := middleware.Logger(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for GET /games/{gameID} to reduce polling noise
		if r.Method == "GET" && regexGetGame.MatchString(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		logger.ServeHTTP(w, r)
	})
}
