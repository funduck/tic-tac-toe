package user

type UserRepo interface {
	FindByID(id string) (*User, error)
	FindBySessionID(sessionID string) (*User, error)
	Save(user *User) error
}
