package app

import (
	"context"

	"github.com/funduck/tic-tac-toe/client/lib"
)

func LoginOrSignup(ctx context.Context, userSvc *lib.UserService, displaySvc *lib.DisplayService) error {
	// Try to login
	token, err := userSvc.Login(ctx, lib.UserID, lib.Password)
	if err == nil {
		lib.AccessToken = *token.AccessToken
		lib.RefreshToken = *token.RefreshToken
		displaySvc.PrintInfo("User logged in successfully.")
		return nil
	}

	// Try to register if login failed
	token, err = userSvc.Signup(ctx, lib.UserID, lib.Password)
	if err != nil {
		return err
	}
	lib.AccessToken = *token.AccessToken
	lib.RefreshToken = *token.RefreshToken
	displaySvc.PrintInfo("User registered successfully.")
	return nil
}

func RefreshToken(ctx context.Context, userSvc *lib.UserService) error {
	token, err := userSvc.RefreshToken(ctx, lib.UserID, lib.RefreshToken)
	if err != nil {
		return err
	}
	lib.AccessToken = *token.AccessToken
	lib.RefreshToken = *token.RefreshToken
	return nil
}
