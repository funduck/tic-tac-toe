package user

type User struct {
	ID           string `json:"id"`
	Password     string
	RefreshToken string
}
