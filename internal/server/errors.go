package server

import (
	"errors"
	"net/http"

	"github.com/funduck/tic-tac-toe/internal/auth"
	"github.com/funduck/tic-tac-toe/internal/game"
	"github.com/funduck/tic-tac-toe/internal/user"
)

type errorEntry struct {
	sentinel error
	status   int
	code     ErrorCode
}

// This registry connects domain errors, http codes and error codes sent to client
var errorRegistry = []errorEntry{
	{game.ErrGameNotFound, http.StatusNotFound, CodeGameNotFound},
	{game.ErrGameNotWaiting, http.StatusConflict, CodeGameNotWaiting},
	{game.ErrGameNotActive, http.StatusConflict, CodeGameNotActive},
	{game.ErrNotYourTurn, http.StatusConflict, CodeNotYourTurn},
	{game.ErrNotInGame, http.StatusConflict, CodeNotInGame},
	{game.ErrCellOccupied, http.StatusBadRequest, CodeCellOccupied},
	{game.ErrOutOfBounds, http.StatusBadRequest, CodeOutOfBounds},
	{user.ErrPasswordIsTooShort, http.StatusBadRequest, CodePasswordTooShort},
	{user.ErrUserAlreadyExists, http.StatusConflict, CodeUserAlreadyExists},
	{user.ErrUserNotFound, http.StatusUnauthorized, CodeUserNotFound},
	{user.ErrInvalidCredentials, http.StatusUnauthorized, CodeInvalidCredentials},
	{user.ErrRefreshTokenDeleted, http.StatusUnauthorized, CodeRefreshTokenDeleted},
	{auth.ErrTokenInvalid, http.StatusUnauthorized, CodeTokenInvalid},
	{auth.ErrTokenExpired, http.StatusUnauthorized, CodeTokenExpired},
	{auth.ErrTokenSignatureInvalid, http.StatusUnauthorized, CodeTokenSignatureInvalid},
}

// lookupError returns the HTTP status and machine-readable code for a domain error.
// Falls back to 500 and empty code for unrecognised errors.
func lookupError(err error) (status int, code ErrorCode) {
	for _, e := range errorRegistry {
		if errors.Is(err, e.sentinel) {
			return e.status, e.code
		}
	}
	return http.StatusInternalServerError, ""
}
