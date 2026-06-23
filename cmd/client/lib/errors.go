package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
)

// APIError is a structured error returned by the server with a machine-readable code.
type APIError struct {
	Code    openapi.ServerErrorCode
	Message string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// HasCode reports whether the error carries the given code.
func (e *APIError) HasCode(code openapi.ServerErrorCode) bool {
	return e.Code == code
}

// friendlyMessages maps server error codes to short, user-facing text.
var friendlyMessages = map[openapi.ServerErrorCode]string{
	openapi.CodeCellOccupied:          "That cell is already taken",
	openapi.CodeNotYourTurn:           "It's not your turn yet",
	openapi.CodeOutOfBounds:           "That cell is off the board",
	openapi.CodeGameNotFound:          "That game no longer exists",
	openapi.CodeGameNotWaiting:        "That game is already full",
	openapi.CodeGameNotActive:         "That game isn't in progress",
	openapi.CodeNotInGame:             "You're not part of that game",
	openapi.CodeInvalidCredentials:    "Wrong username or password",
	openapi.CodeUserNotFound:          "Wrong username or password",
	openapi.CodeUserAlreadyExists:     "That username is taken",
	openapi.CodePasswordTooShort:      "Password is too short (minimum 6 characters)",
	openapi.CodeRefreshTokenDeleted:   "Session expired, please restart",
	openapi.CodeTokenInvalid:          "Session expired, please restart",
	openapi.CodeTokenExpired:          "Session expired, please restart",
	openapi.CodeTokenSignatureInvalid: "Session expired, please restart",
}

// FriendlyMessage converts an error into short, user-facing text. Known API
// error codes are mapped to plain language, network failures get a hint about
// the server, and anything else falls back to the underlying message.
func FriendlyMessage(err error) string {
	if err == nil {
		return ""
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		if msg, ok := friendlyMessages[apiErr.Code]; ok {
			return msg
		}
		if apiErr.Message != "" {
			return apiErr.Message
		}
	}

	if isConnectionError(err) {
		return "Can't reach the server — is it running on :8080?"
	}

	return err.Error()
}

// isConnectionError reports whether err looks like a failure to reach the server.
func isConnectionError(err error) bool {
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "dial tcp")
}

// parseError decodes a server JSON error response into an *APIError.
// Falls back to the original err when the response body is absent or unparseable.
func parseError(r *http.Response, err error) error {
	if r == nil {
		return err
	}
	var e openapi.ServerErrorResponse
	if decodeErr := json.NewDecoder(r.Body).Decode(&e); decodeErr != nil {
		return err
	}
	return &APIError{
		Code:    e.GetCode(),
		Message: e.GetError(),
	}
}
