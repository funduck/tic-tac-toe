package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

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
