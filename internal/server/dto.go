package server

type UserSignupRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// AccessTokenResponse is the body returned by the auth endpoints. The refresh
// token is delivered separately as an HttpOnly cookie, never in the body.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// ErrorCode is the machine-readable identifier returned in every error response.
// Should simplify error handling on the client.
type ErrorCode string

const (
	CodeGameNotFound          ErrorCode = "ERR_GAME_NOT_FOUND"
	CodeGameNotWaiting        ErrorCode = "ERR_GAME_NOT_WAITING"
	CodeGameNotActive         ErrorCode = "ERR_GAME_NOT_ACTIVE"
	CodeNotYourTurn           ErrorCode = "ERR_NOT_YOUR_TURN"
	CodeNotInGame             ErrorCode = "ERR_NOT_IN_GAME"
	CodeCellOccupied          ErrorCode = "ERR_CELL_OCCUPIED"
	CodeOutOfBounds           ErrorCode = "ERR_OUT_OF_BOUNDS"
	CodePasswordTooShort      ErrorCode = "ERR_PASSWORD_TOO_SHORT"
	CodeUserAlreadyExists     ErrorCode = "ERR_USER_ALREADY_EXISTS"
	CodeUserNotFound          ErrorCode = "ERR_USER_NOT_FOUND"
	CodeInvalidCredentials    ErrorCode = "ERR_INVALID_CREDENTIALS"
	CodeRefreshTokenDeleted   ErrorCode = "ERR_REFRESH_TOKEN_DELETED"
	CodeTokenInvalid          ErrorCode = "ERR_TOKEN_INVALID"
	CodeTokenExpired          ErrorCode = "ERR_TOKEN_EXPIRED"
	CodeTokenSignatureInvalid ErrorCode = "ERR_TOKEN_SIGNATURE_INVALID"
)

type ErrorResponse struct {
	Error string    `json:"error"`
	Code  ErrorCode `json:"code,omitempty"`
}
