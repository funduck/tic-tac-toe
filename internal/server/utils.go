package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

const (
	// refreshTokenCookieName is the cookie carrying the refresh token. It is kept
	// out of the JSON body so a browser client never exposes it to JavaScript (XSS).
	refreshTokenCookieName = "refresh_token"
	// refreshTokenCookiePath scopes the cookie so browsers only send it to the
	// refresh endpoint.
	refreshTokenCookiePath = "/api/users/refresh-token"
	// refreshTokenCookieMaxAge mirrors the refresh token lifetime in
	// auth.AccessTokenService (7 days).
	refreshTokenCookieMaxAge = 7 * 24 * 60 * 60
)

// setRefreshTokenCookie writes the refresh token as an HttpOnly cookie. secure is
// driven by config so the cookie also works over plain http in local development.
func setRefreshTokenCookie(w http.ResponseWriter, token string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    token,
		Path:     refreshTokenCookiePath,
		MaxAge:   refreshTokenCookieMaxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}

// readRefreshTokenCookie returns the refresh token from the request cookie.
func readRefreshTokenCookie(r *http.Request) (string, error) {
	c, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func parseRequestBody(r *http.Request, dst any, logger *slog.Logger, w http.ResponseWriter) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // catch unknown fields
	if err := decoder.Decode(dst); err != nil {
		if logger != nil {
			logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err)
		}
		writeError(w, http.StatusBadRequest, err)
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// writeDomainError writes a JSON error response with status and code both resolved from the registry.
func writeDomainError(w http.ResponseWriter, err error) {
	status, code := lookupError(err)
	writeJSON(w, status, ErrorResponse{Error: err.Error(), Code: code})
}

// writeError writes a JSON error response with a caller-supplied status; code is still looked up and attached.
func writeError(w http.ResponseWriter, status int, err error) {
	_, code := lookupError(err)
	writeJSON(w, status, ErrorResponse{Error: err.Error(), Code: code})
}

func userIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDContextKey).(string)
	return v
}
