package server

import (
	"errors"
	"net/http"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/game"
	"github.com/funduck/tic-tac-toe/internal/user"
)

func TestLookupError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   ErrorCode
	}{
		{game.ErrGameNotFound.Error(), game.ErrGameNotFound, http.StatusNotFound, CodeGameNotFound},
		{game.ErrGameNotWaiting.Error(), game.ErrGameNotWaiting, http.StatusConflict, CodeGameNotWaiting},
		{game.ErrGameNotActive.Error(), game.ErrGameNotActive, http.StatusConflict, CodeGameNotActive},
		{game.ErrNotYourTurn.Error(), game.ErrNotYourTurn, http.StatusConflict, CodeNotYourTurn},
		{game.ErrNotInGame.Error(), game.ErrNotInGame, http.StatusConflict, CodeNotInGame},
		{game.ErrCellOccupied.Error(), game.ErrCellOccupied, http.StatusBadRequest, CodeCellOccupied},
		{game.ErrOutOfBounds.Error(), game.ErrOutOfBounds, http.StatusBadRequest, CodeOutOfBounds},
		{user.ErrPasswordIsTooShort.Error(), user.ErrPasswordIsTooShort, http.StatusBadRequest, CodePasswordTooShort},
		{user.ErrUserAlreadyExists.Error(), user.ErrUserAlreadyExists, http.StatusConflict, CodeUserAlreadyExists},
		{user.ErrUserNotFound.Error(), user.ErrUserNotFound, http.StatusUnauthorized, CodeUserNotFound},
		{user.ErrInvalidCredentials.Error(), user.ErrInvalidCredentials, http.StatusUnauthorized, CodeInvalidCredentials},
		{user.ErrRefreshTokenDeleted.Error(), user.ErrRefreshTokenDeleted, http.StatusUnauthorized, CodeRefreshTokenDeleted},
		{auth.ErrTokenInvalid.Error(), auth.ErrTokenInvalid, http.StatusUnauthorized, CodeTokenInvalid},
		{auth.ErrTokenExpired.Error(), auth.ErrTokenExpired, http.StatusUnauthorized, CodeTokenExpired},
		{auth.ErrTokenSignatureInvalid.Error(), auth.ErrTokenSignatureInvalid, http.StatusUnauthorized, CodeTokenSignatureInvalid},
		{"unknown error", errors.New("some unexpected error"), http.StatusInternalServerError, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, code := lookupError(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}
			if code != tt.expectedCode {
				t.Errorf("expected code %q, got %q", tt.expectedCode, code)
			}
		})
	}
}
