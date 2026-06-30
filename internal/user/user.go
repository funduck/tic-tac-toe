package user

type User struct {
	ID           string `json:"id"`
	Password     string // hashed password
	RefreshToken string // hashed refresh token
}

const defaultPasswordLength = 6
