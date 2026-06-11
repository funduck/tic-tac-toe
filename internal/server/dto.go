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

type ErrorResponse struct {
	Error string `json:"error"`
}
