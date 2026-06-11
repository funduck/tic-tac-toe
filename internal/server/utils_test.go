package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/game"
)

func TestGameErrorMapping(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "game not found",
			err:      game.ErrGameNotFound,
			expected: http.StatusNotFound,
		},
		{
			name:     "game not waiting",
			err:      game.ErrGameNotWaiting,
			expected: http.StatusConflict,
		},
		{
			name:     "game not active",
			err:      game.ErrGameNotActive,
			expected: http.StatusConflict,
		},
		{
			name:     "not in game",
			err:      game.ErrNotInGame,
			expected: http.StatusConflict,
		},
		{
			name:     "not your turn",
			err:      game.ErrNotYourTurn,
			expected: http.StatusConflict,
		},
		{
			name:     "cell occupied",
			err:      game.ErrCellOccupied,
			expected: http.StatusBadRequest,
		},
		{
			name:     "unknown error",
			err:      errors.New("some unexpected error"),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := mapGameError(tt.err)
			if status != tt.expected {
				t.Errorf("expected status %d, got %d", tt.expected, status)
			}
		})
	}
}

func TestParseRequestBody(t *testing.T) {
	type testRequest struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	tests := []struct {
		name        string
		body        interface{}
		expected    testRequest
		expectError bool
	}{
		{
			name:        "valid JSON",
			body:        testRequest{Field1: "test", Field2: 42},
			expected:    testRequest{Field1: "test", Field2: 42},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			body:        "not a json",
			expected:    testRequest{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(bodyBytes))
			var parsed testRequest
			err := parseRequestBody(req, &parsed, slog.New(slog.NewTextHandler(io.Discard, nil)), httptest.NewRecorder())
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && parsed != tt.expected {
				t.Errorf("expected %+v, got %+v", tt.expected, parsed)
			}
		})
	}
}
