package user

import "context"

type UserRepo interface {
	FindByID(ctx context.Context, id string) (*User, error)
	// Create inserts a new user or fails with ErrUserAlreadyExists.
	Create(ctx context.Context, user *User) error
	// Save updates an existing user (or inserts unconditionally).
	Save(ctx context.Context, user *User) error
}
