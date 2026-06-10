package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func parseRequestBody(r *http.Request, dst any, logger *slog.Logger, w http.ResponseWriter) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // catch unknown fields
	if err := decoder.Decode(dst); err != nil {
		logger.Warn("bad request", "method", r.Method, "path", r.URL.Path, "error", err)
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

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrorResponse{Error: err.Error()})
}
