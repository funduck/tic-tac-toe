package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

type TokenService interface {
	GenerateTokens(userID string) (*auth.TokenPair, error)
	ValidateToken(tokenStr string) (*auth.Token, error)
	ValidateRefreshToken(tokenStr string) (string, error)
}

var (
	ErrPasswordIsTooShort  = errors.New("password is too short")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrRefreshTokenDeleted = errors.New("refresh token deleted")
)

// UserService provides user management and authentication services using
// a UserRepo for persistence and a TokenService for token generation and validation.
type UserService struct {
	userRepo     UserRepo
	tokenService TokenService
	logger       *slog.Logger
}

func NewUserService(userRepo UserRepo, tokenService TokenService, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenService: tokenService,
		logger:       logger,
	}
}

// Signup creates a new user and returns a token pair.
func (s *UserService) Signup(ctx context.Context, userID, password string) (*User, *auth.TokenPair, error) {
	if len(password) < defaultPasswordLength {
		return nil, nil, fmt.Errorf("signup: %w, should contain at least %d characters", ErrPasswordIsTooShort, defaultPasswordLength)
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	tokenPair, err := s.tokenService.GenerateTokens(userID)
	if err != nil {
		return nil, nil, err
	}
	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user := &User{
		ID:           userID,
		Password:     passwordHash,
		RefreshToken: refreshTokenHash,
	}
	// Instead of Save we use Create to ensure that we don't overwrite an existing user.
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("user signed up", "userID", userID)

	return user, tokenPair, nil
}

// Login authenticates a user and returns a token pair.
func (s *UserService) Login(ctx context.Context, userID, password string) (*User, *auth.TokenPair, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, nil, ErrInvalidCredentials
	}

	tokenPair, err := s.tokenService.GenerateTokens(userID)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("user logged in", "userID", userID)

	return user, tokenPair, nil
}

// RefreshToken validates the refresh token and returns a new token pair if valid.
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (*User, *auth.TokenPair, error) {
	userID, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if !auth.ValidateRefreshTokenHash(refreshToken, user.RefreshToken) {
		return nil, nil, ErrRefreshTokenDeleted
	}

	tokenPair, err := s.tokenService.GenerateTokens(userID)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("user refreshed token", "userID", userID)

	return user, tokenPair, nil
}
