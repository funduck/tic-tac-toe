package server

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

type TokenService interface {
	ValidateToken(tokenStr string) (token *auth.Token, err error)
}

var reAuthorization = regexp.MustCompile(`^Bearer (.+)$`)

type contextKey string

const userIDContextKey contextKey = "userID"

// AuthMiddleware is a middleware that handles authentication
func AuthMiddleware(tokenService TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract authorization header and validate JWT token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			// Validate token format (e.g., "Bearer <token>")
			tokenStr := ""
			matches := reAuthorization.FindStringSubmatch(authHeader)
			if len(matches) == 2 {
				tokenStr = matches[1]
			} else {
				writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))
				return
			}

			// Validate JWT token
			token, err := tokenService.ValidateToken(tokenStr)
			if err != nil {
				writeDomainError(w, err)
				return
			}
			userID := token.UserID

			// Propagate userID to the request context for downstream handlers
			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			r = r.WithContext(ctx)

			// If valid, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
