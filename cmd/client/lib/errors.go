package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

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
