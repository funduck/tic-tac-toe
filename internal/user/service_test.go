package user

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

type MockUserRepo struct {
	users map[string]*User
}

func NewMockUserRepo() *MockUserRepo {
	return &MockUserRepo{
		users: make(map[string]*User),
	}
}

func (r *MockUserRepo) FindByID(id string) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *MockUserRepo) FindBySessionID(sessionID string) (*User, error) {
	for _, user := range r.users {
		if user.ID == sessionID { // This is just a placeholder. In a real implementation, you'd check the session ID.
			return user, nil
		}
	}
	return nil, nil
}

func (r *MockUserRepo) Save(user *User) error {
	r.users[user.ID] = user
	return nil
}

func TestAuthService_Signup(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	userID := "testuser"
	password := "password123"

	user, tokenPair, err := authService.Signup(userID, password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected non-nil user")
	}
	if user.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, user.ID)
	}
	if tokenPair == nil {
		t.Fatal("expected non-nil token pair")
	}
	if tokenPair.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestAuthService_Signup_Errors(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	// prepare existing user
	_, _, err := authService.Signup("testuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error during initial signup: %v", err)
	}

	for _, tc := range []struct {
		name          string
		userID        string
		password      string
		expectedError error
	}{
		{"password too short", "testuser1", "123", ErrPasswordIsTooShort},
		{"user already exists", "testuser", "password123", ErrUserAlreadyExists},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := authService.Signup(tc.userID, tc.password)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	// First, sign up a user
	userID := "testuser"
	password := "password123"
	user, _, err := authService.Signup(userID, password)
	if err != nil {
		t.Fatalf("unexpected error during signup: %v", err)
	}

	// Now, try to log in with the same credentials
	loggedInUser, tokenPair, err := authService.Login(userID, password)
	if err != nil {
		t.Fatalf("unexpected error during login: %v", err)
	}
	if loggedInUser == nil {
		t.Fatal("expected non-nil user on login")
	}
	if loggedInUser.ID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, loggedInUser.ID)
	}
	if tokenPair == nil {
		t.Fatal("expected non-nil token pair on login")
	}
	if tokenPair.AccessToken == "" {
		t.Error("expected non-empty access token on login")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("expected non-empty refresh token on login")
	}
}

func TestAuthService_Login_Errors(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	// prepare existing user
	_, _, err := authService.Signup("testuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error during signup: %v", err)
	}

	for _, tc := range []struct {
		name          string
		userID        string
		password      string
		expectedError error
	}{
		{"nonexistent user", "nonexistent", "password123", ErrUserNotFound},
		{"wrong password", "testuser", "wrongpassword", ErrInvalidCredentials},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := authService.Login(tc.userID, tc.password)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	// First, sign up a user
	userID := "testuser"
	password := "password123"
	user, tokenPair, err := authService.Signup(userID, password)
	if err != nil {
		t.Fatalf("unexpected error during signup: %v", err)
	}

	// Now, try to refresh the token using the refresh token
	refreshedUser, newTokenPair, err := authService.RefreshToken(userID, tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error during token refresh: %v", err)
	}
	if refreshedUser == nil {
		t.Fatal("expected non-nil user on token refresh")
	}
	if refreshedUser.ID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, refreshedUser.ID)
	}
	if newTokenPair == nil {
		t.Fatal("expected non-nil token pair on token refresh")
	}
	if newTokenPair.AccessToken == "" {
		t.Error("expected non-empty access token on token refresh")
	}
	if newTokenPair.RefreshToken == "" {
		t.Error("expected non-empty refresh token on token refresh")
	}
}

func TestAuthService_RefreshToken_Errors(t *testing.T) {
	userRepo := NewMockUserRepo()
	tokenService := auth.NewAccessTokenService("secret", "my-awesome-app")
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	authService := NewUserService(userRepo, tokenService, logger)

	// prepare valid refresh token
	_, tokenPair, err := authService.Signup("testuser", "password123")
	if err != nil {
		t.Fatalf("unexpected error during signup: %v", err)
	}
	userRepo.users["testuser"] = nil // Simulate user not found for the refresh token test

	// prepare existing user but refresh token hash will not match
	_, tokenPair2, err := authService.Signup("testuser2", "password123")
	if err != nil {
		t.Fatalf("unexpected error during signup: %v", err)
	}
	userRepo.users["testuser2"].RefreshToken = "invalidhash" // Simulate refresh token hash mismatch

	for _, tc := range []struct {
		name          string
		userID        string
		refreshToken  string
		expectedError error
	}{
		{"invalid refresh token", "testuser", "invalidtoken", auth.ErrTokenInvalid},
		{"token user not matches user", "nonexistent", tokenPair.RefreshToken, ErrUserNotFound},
		{"user not found", "testuser", tokenPair.RefreshToken, ErrUserNotFound},
		{"refresh token hash mismatch", "testuser2", tokenPair2.RefreshToken, ErrRefreshTokenDeleted},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := authService.RefreshToken(tc.userID, tc.refreshToken)
			if !errors.Is(err, tc.expectedError) {
				t.Fatalf("expected error %v, got %v", tc.expectedError, err)
			}
		})
	}
}
