package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

type TokenService interface {
	ValidateToken(tokenStr string) (token *auth.Token, err error)
}

var reAuthorization = regexp.MustCompile(`^Bearer (.+)$`)

type dtoWithUserID struct {
	UserID string `json:"userID"`
}

// AuthMiddleware is a middleware that handles authentication
func AuthMiddleware(tokenService TokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract authorization header and validate JWT token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate token format (e.g., "Bearer <token>")
			tokenStr := ""
			matches := reAuthorization.FindStringSubmatch(authHeader)
			if len(matches) == 2 {
				tokenStr = matches[1]
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate JWT token
			token, err := tokenService.ValidateToken(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			userID := token.UserID

			// check if body contains userID and if so, validate it matches the token's userID
			if r.Body != nil && r.Header.Get("Content-Type") == "application/json" {
				bodyBytes, err := getBodyBytes(r)
				if err != nil {
					http.Error(w, "Failed to read body", http.StatusBadRequest)
					return
				}

				var dto dtoWithUserID
				if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&dto); err == nil {
					if dto.UserID != userID {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
						return
					}
				}
			}

			// If valid, call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
