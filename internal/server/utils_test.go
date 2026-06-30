package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
