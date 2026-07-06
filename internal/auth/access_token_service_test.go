package auth

import (
	"errors"
	"testing"
	"time"
)

func TestAccessTokenService_GenerateTokens(t *testing.T) {
	s := NewAccessTokenService(
		"secret",
		"my-awesome-app",
	)

	tokenPair, err := s.GenerateTokens("user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenPair == nil {
		t.Fatal("expected non-nil token")
	}
	if tokenPair.AccessToken == "" {
		t.Error("expected non-nil access token")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("expected non-nil refresh token")
	}
}

func TestAccessTokenService_ValidateToken(t *testing.T) {
	s := NewAccessTokenService(
		"secret",
		"my-awesome-app",
	)

	tokenPair, err := s.GenerateTokens("user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	token, err := s.ValidateToken(tokenPair.AccessToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == nil {
		t.Fatal("expected non-nil token")
	}
	if token.Subject != "user1" {
		t.Errorf("expected user ID 'user1', got '%s'", token.Subject)
	}
}

func TestAccessTokenService_ValidateToken_Errors(t *testing.T) {
	s := NewAccessTokenService(
		"secret",
		"my-awesome-app",
	)

	for _, tc := range []struct {
		name        string
		tokenStr    string
		expectedErr error
	}{
		{
			name:        "invalid token format",
			tokenStr:    "invalid-token",
			expectedErr: ErrTokenInvalid,
		},
		{
			name: "token signed with different secret",
			tokenStr: func() string {
				s2 := NewAccessTokenService(
					"wrong-secret",
					"my-awesome-app",
				)
				tokenPair, err := s2.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.AccessToken
			}(),
			expectedErr: ErrTokenInvalid,
		},
		{
			name: "expired token",
			tokenStr: func() string {
				// Create a token with a short lifetime for testing
				s2 := &AccessTokenService{
					secret:              "secret",
					issuer:              "my-awesome-app",
					accessTokenLifetime: -1 * time.Minute, // Already expired
				}
				tokenPair, err := s2.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.AccessToken
			}(),
			expectedErr: ErrTokenExpired,
		},
		{
			name: "refresh token used as access token",
			tokenStr: func() string {
				tokenPair, err := s.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.RefreshToken
			}(),
			expectedErr: ErrTokenInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.ValidateToken(tc.tokenStr)
			if err == nil {
				t.Fatal("expected error for invalid token")
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}

}

func TestAccessTokenService_ValidateRefreshToken(t *testing.T) {
	s := NewAccessTokenService(
		"secret",
		"my-awesome-app",
	)

	tokenPair, err := s.GenerateTokens("user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	userID, err := s.ValidateRefreshToken(tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user1" {
		t.Errorf("expected user ID 'user1', got '%s'", userID)
	}
}

func TestAccessTokenService_ValidateRefreshToken_Errors(t *testing.T) {
	s := NewAccessTokenService(
		"secret",
		"my-awesome-app",
	)

	for _, tc := range []struct {
		name        string
		tokenStr    string
		expectedErr error
	}{
		{
			name:        "invalid token format",
			tokenStr:    "invalid-token",
			expectedErr: ErrTokenInvalid,
		},
		{
			name: "token signed with different secret",
			tokenStr: func() string {
				s2 := NewAccessTokenService(
					"wrong-secret",
					"my-awesome-app",
				)
				tokenPair, err := s2.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.RefreshToken
			}(),
			expectedErr: ErrTokenInvalid,
		},
		{
			name: "expired token",
			tokenStr: func() string {
				// Create a token with a short lifetime for testing
				s2 := &AccessTokenService{
					secret:               "secret",
					issuer:               "my-awesome-app",
					refreshTokenLifetime: -1 * time.Minute, // Already expired
				}
				tokenPair, err := s2.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.RefreshToken
			}(),
			expectedErr: ErrTokenInvalid, // For refresh token expiration means same as invalid token.
		},
		{
			name: "access token used as refresh token",
			tokenStr: func() string {
				tokenPair, err := s.GenerateTokens("user1")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return tokenPair.AccessToken
			}(),
			expectedErr: ErrTokenInvalid,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := s.ValidateRefreshToken(tc.tokenStr)
			if err == nil {
				t.Fatal("expected error for invalid token")
			}
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
