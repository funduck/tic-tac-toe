package lib

import (
	"errors"
	"fmt"
	"net"
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
