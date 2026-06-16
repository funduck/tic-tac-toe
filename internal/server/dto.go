package server

type CreateGameRequest struct {
	GameID  string `json:"gameID"`
	UserID  string `json:"userID"`
	Private bool   `json:"private"`
}

type JoinGameRequest struct {
	GameID string `json:"gameID"`
	UserID string `json:"userID"`
}

type JoinAnyGameRequest struct {
	UserID string `json:"userID"`
}

type MoveRequest struct {
	GameID string `json:"gameID"`
	UserID string `json:"userID"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type GiveUpRequest struct {
	GameID string `json:"gameID"`
	UserID string `json:"userID"`
}

type UserSignupRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type UserLoginRequest struct {
	UserID   string `json:"userID"`
	Password string `json:"password"`
}

type UserRefreshTokenRequest struct {
	UserID       string `json:"userID"`
	RefreshToken string `json:"refreshToken"`
}

// ErrorCode is the machine-readable identifier returned in every error response.
type ErrorCode string

const (
	CodeGameNotFound           ErrorCode = "ERR_GAME_NOT_FOUND"
	CodeGameNotWaiting         ErrorCode = "ERR_GAME_NOT_WAITING"
	CodeGameNotActive          ErrorCode = "ERR_GAME_NOT_ACTIVE"
	CodeNotYourTurn            ErrorCode = "ERR_NOT_YOUR_TURN"
	CodeNotInGame              ErrorCode = "ERR_NOT_IN_GAME"
	CodeCellOccupied           ErrorCode = "ERR_CELL_OCCUPIED"
	CodeOutOfBounds            ErrorCode = "ERR_OUT_OF_BOUNDS"
	CodePasswordTooShort       ErrorCode = "ERR_PASSWORD_TOO_SHORT"
	CodeUserAlreadyExists      ErrorCode = "ERR_USER_ALREADY_EXISTS"
	CodeUserNotFound           ErrorCode = "ERR_USER_NOT_FOUND"
	CodeInvalidCredentials     ErrorCode = "ERR_INVALID_CREDENTIALS"
	CodeRefreshTokenDeleted    ErrorCode = "ERR_REFRESH_TOKEN_DELETED"
	CodeTokenInvalid           ErrorCode = "ERR_TOKEN_INVALID"
	CodeTokenExpired           ErrorCode = "ERR_TOKEN_EXPIRED"
	CodeTokenSignatureInvalid  ErrorCode = "ERR_TOKEN_SIGNATURE_INVALID"
)

type ErrorResponse struct {
	Error string    `json:"error"`
	Code  ErrorCode `json:"code,omitempty"`
}
