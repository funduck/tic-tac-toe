package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessTokenService struct {
	secret               string
	issuer               string
	accessTokenLifetime  time.Duration
	refreshTokenLifetime time.Duration
}

func NewAccessTokenService(secret, issuer string) *AccessTokenService {
	return &AccessTokenService{
		secret:               secret,
		issuer:               issuer,
		accessTokenLifetime:  1 * time.Minute,
		refreshTokenLifetime: 7 * 24 * time.Hour,
	}
}

func (s *AccessTokenService) keyFunc(token *jwt.Token) (interface{}, error) {
	// Check the signing method
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrTokenSignatureInvalid
	}
	return []byte(s.secret), nil
}

// Parses and validates the access token, returning the user ID and session ID if valid.
func (s *AccessTokenService) ValidateToken(tokenStr string) (*Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, s.keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, errors.Join(err, ErrTokenInvalid)
	}

	// Check if token is valid and claims are of type *CustomClaims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return &Token{*claims}, nil
	}

	// Type assertion failed or token is invalid
	return nil, ErrTokenInvalid
}

// Parses and validates the refresh token, returning the user ID.
func (s *AccessTokenService) ValidateRefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, s.keyFunc)
	if err != nil {
		// For refresh token expiration means same as invalid token.
		// if errors.Is(err, jwt.ErrTokenExpired) {
		// 	return "", ErrTokenExpired
		// }
		return "", errors.Join(err, ErrTokenInvalid)
	}

	// Check if token is valid and claims are of type *jwt.RegisteredClaims
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}

	// Type assertion failed or token is invalid
	return "", ErrTokenInvalid
}

// Generates JWT token and refresh token for the given user ID. The refresh token is stored in the user repo for later validation.
func (s *AccessTokenService) GenerateToken(userID string) (*TokenPair, error) {
	now := time.Now()

	id, err := uuid.NewV7() // V7 for lexicographically sortable IDs
	if err != nil {
		return nil, err
	}
	sessionID := id.String()

	accessClaims := CustomClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
		},
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenLifetime)),
		IssuedAt:  jwt.NewNumericDate(now),
		Issuer:    s.issuer,
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
