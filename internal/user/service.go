package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/funduck/tic-tac-toe/internal/auth"
)

type TokenService interface {
	GenerateToken(userID string) (*auth.TokenPair, error)
	ValidateToken(tokenStr string) (*auth.Token, error)
	ValidateRefreshToken(tokenStr string) (string, error)
}

var ErrPasswordIsTooShort = errors.New("password is too short")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrRefreshTokenDeleted = errors.New("refresh token deleted")

type UserService struct {
	userRepo     UserRepo
	tokenService TokenService
	secret       string
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
	if len(password) < 6 {
		return nil, nil, fmt.Errorf("signup: %w, should contain at least 6 characters", ErrPasswordIsTooShort)
	}

	existingUser, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, ErrUserAlreadyExists
	}

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}
	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user := &User{
		ID:           userID,
		Password:     passwordHash,
		RefreshToken: refreshTokenHash,
	}
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("User signed up", "userID", userID)

	return user, tokenPair, nil
}

// Login authenticates a user and returns a token pair.
func (s *UserService) Login(ctx context.Context, userID, password string) (*User, *auth.TokenPair, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, errors.Join(err, ErrUserNotFound)
	}

	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, nil, ErrInvalidCredentials
	}

	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("User logged in", "userID", userID)

	return user, tokenPair, nil
}

// RefreshToken validates the refresh token and returns a new token pair if valid.
func (s *UserService) RefreshToken(ctx context.Context, userID, refreshToken string) (*User, *auth.TokenPair, error) {
	subject, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	if userID != subject {
		return nil, nil, ErrUserNotFound
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

	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, nil, err
	}

	s.logger.Info("User refreshed token", "userID", userID)

	return user, tokenPair, nil
}
