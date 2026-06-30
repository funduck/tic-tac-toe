package user

type User struct {
	ID           string `json:"id"`
	Password     string // hashed password
	RefreshToken string // hashed refresh token
}

func (u *User) Clone() *User {
	return &User{
		ID:           u.ID,
		Password:     u.Password,
		RefreshToken: u.RefreshToken,
	}
}

const defaultPasswordLength = 6
