package auth

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrTokenInvalid = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")
var ErrTokenSignatureInvalid = errors.New("invalid token signature")

type CustomClaims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	CustomClaims
}

// HashRefreshToken takes a long JWT string and returns a fast, fixed-length secure hash
func HashRefreshToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

// ValidateRefreshTokenHash hashes the incoming token and safely compares it with the DB hash
func ValidateRefreshTokenHash(incomingToken string, dbHash string) bool {
	// 1. Calculate the SHA-256 hash of the token sent by the client
	hasher := sha256.New()
	hasher.Write([]byte(incomingToken))
	computedHashBytes := hasher.Sum(nil)

	// 2. Decode the hex hash from your database back into bytes
	dbHashBytes, err := hex.DecodeString(dbHash)
	if err != nil {
		return false
	}

	// 3. SECURE COMPARISON: ConstantTimeCompare prevents timing attacks.
	// It returns 1 if they match, and 0 if they don't.
	if subtle.ConstantTimeCompare(computedHashBytes, dbHashBytes) == 1 {
		return true
	}

	return false
}
