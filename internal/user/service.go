package user

import (
	"errors"

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
}

func NewUserService(userRepo UserRepo, tokenService TokenService) *UserService {
	return &UserService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Signup creates a new user and returns a token pair.
func (s *UserService) Signup(userID, password string) (*User, *auth.TokenPair, error) {
	// Check password length
	if len(password) < 6 {
		return nil, nil, ErrPasswordIsTooShort
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, ErrUserAlreadyExists
	}

	// Generate password hash
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, nil, err
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}
	// Hash the refresh token before storing it
	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user := &User{
		ID:           userID,
		Password:     passwordHash,
		RefreshToken: refreshTokenHash,
	}
	if err := s.userRepo.Save(user); err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

// Login authenticates a user and returns a token pair.
func (s *UserService) Login(userID, password string) (*User, *auth.TokenPair, error) {
	// Find user by ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, errors.Join(err, ErrUserNotFound)
	}

	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, nil, ErrInvalidCredentials
	}

	// Generate token pair
	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}

	// Hash the refresh token before storing it
	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(user); err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

// RefreshToken validates the refresh token and returns a new token pair if valid.
func (s *UserService) RefreshToken(userID, refreshToken string) (*User, *auth.TokenPair, error) {
	// Validate refresh token
	subject, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, err
	}

	if userID != subject {
		return nil, nil, ErrUserNotFound
	}

	// Find user by ID
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, ErrUserNotFound
	}

	// Validate the refresh token against the stored hash
	if !auth.ValidateRefreshTokenHash(refreshToken, user.RefreshToken) {
		return nil, nil, ErrRefreshTokenDeleted
	}

	// Generate new token pair
	tokenPair, err := s.tokenService.GenerateToken(userID)
	if err != nil {
		return nil, nil, err
	}

	// Hash the refresh token before storing it
	refreshTokenHash := auth.HashRefreshToken(tokenPair.RefreshToken)

	user.RefreshToken = refreshTokenHash
	if err := s.userRepo.Save(user); err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}
