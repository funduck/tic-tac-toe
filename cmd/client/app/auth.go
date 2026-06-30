package app

import (
	"context"
	"errors"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID"
	"github.com/funduck/tic-tac-toe/client/lib"
)

func LoginOrSignup(ctx context.Context, userSvc *lib.UserService, displaySvc *lib.DisplayService) error {
	// Try to login
	token, err := userSvc.Login(ctx, lib.UserID, lib.Password)
	if err == nil {
		lib.AccessToken = *token.AccessToken
		displaySvc.PrintInfo("User logged in successfully.")
		return nil
	}

	// Only fall through to registration if the "user not found" error occurred
	var apiErr *lib.APIError
	if !errors.As(err, &apiErr) || !apiErr.HasCode(openapi.CodeUserNotFound) {
		return err
	}

	// Try to register if login failed
	displaySvc.PrintInfo("Attempting to register a new user...")
	token, err = userSvc.Signup(ctx, lib.UserID, lib.Password)
	if err != nil {
		return err
	}
	lib.AccessToken = *token.AccessToken
	displaySvc.PrintInfo("User registered successfully.")
	return nil
}

func RefreshToken(ctx context.Context, userSvc *lib.UserService) error {
	token, err := userSvc.RefreshToken(ctx)
	if err != nil {
		return err
	}
	lib.AccessToken = *token.AccessToken
	return nil
}
