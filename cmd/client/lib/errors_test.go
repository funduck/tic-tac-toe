package lib

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"testing"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func TestFriendlyMessage(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
		{
			name: "known code is mapped",
			err:  &APIError{Code: openapi.CodeCellOccupied, Message: "ERR_CELL_OCCUPIED"},
			want: "That cell is already taken",
		},
		{
			name: "auth code is mapped",
			err:  &APIError{Code: openapi.CodeInvalidCredentials},
			want: "Wrong username or password",
		},
		{
			name: "unknown code falls back to message",
			err:  &APIError{Code: "ERR_SOMETHING_NEW", Message: "raw server message"},
			want: "raw server message",
		},
		{
			name: "connection refused string",
			err:  errors.New("dial tcp 127.0.0.1:8080: connect: connection refused"),
			want: "Can't reach the server — is it running on :8080?",
		},
		{
			name: "net.OpError",
			err:  &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("refused")},
			want: "Can't reach the server — is it running on :8080?",
		},
		{
			name: "refresh token deleted maps to session message",
			err:  &APIError{Code: openapi.CodeRefreshTokenDeleted},
			want: sessionEndedMessage,
		},
		{
			name: "wrapped API error",
			err:  fmt.Errorf("make move: %w", &APIError{Code: openapi.CodeNotYourTurn}),
			want: "It's not your turn yet",
		},
		{
			name: "generic error passthrough",
			err:  errors.New("something odd"),
			want: "something odd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FriendlyMessage(tt.err); got != tt.want {
				t.Fatalf("FriendlyMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsSessionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"refresh token deleted", &APIError{Code: openapi.CodeRefreshTokenDeleted}, true},
		{"token expired", &APIError{Code: openapi.CodeTokenExpired}, true},
		{"token invalid", &APIError{Code: openapi.CodeTokenInvalid}, true},
		{"token signature invalid", &APIError{Code: openapi.CodeTokenSignatureInvalid}, true},
		{"unrelated api error", &APIError{Code: openapi.CodeNotYourTurn}, false},
		{"plain error", errors.New("boom"), false},
		{
			// The HTTP client wraps a RoundTrip error in *url.Error; IsSessionError
			// must still see through to the underlying *APIError.
			name: "wrapped in url.Error",
			err:  &url.Error{Op: "Get", URL: "http://x", Err: &APIError{Code: openapi.CodeRefreshTokenDeleted}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSessionError(tt.err); got != tt.want {
				t.Fatalf("IsSessionError() = %v, want %v", got, tt.want)
			}
		})
	}
}
