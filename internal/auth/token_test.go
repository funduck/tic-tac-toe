package auth

import "testing"

func TestHashRefreshToken(t *testing.T) {
	hash := HashRefreshToken("sometoken")
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if hash == "sometoken" {
		t.Error("hash must not equal the original token")
	}
}

func TestHashRefreshToken_Deterministic(t *testing.T) {
	if HashRefreshToken("abc") != HashRefreshToken("abc") {
		t.Error("expected same input to always produce the same hash")
	}
}

func TestHashRefreshToken_Distinct(t *testing.T) {
	if HashRefreshToken("token-a") == HashRefreshToken("token-b") {
		t.Error("expected different tokens to produce different hashes")
	}
}

func TestValidateRefreshTokenHash(t *testing.T) {
	token := "some.jwt.token"
	hash := HashRefreshToken(token)
	if !ValidateRefreshTokenHash(token, hash) {
		t.Error("expected token to match its own hash")
	}
}

func TestValidateRefreshTokenHash_Invalid(t *testing.T) {
	token := "some.jwt.token"
	hash := HashRefreshToken(token)

	for _, tc := range []struct {
		name     string
		incoming string
		dbHash   string
	}{
		{"wrong token", "other.jwt.token", hash},
		{"malformed hash", token, "not-valid-hex!!"},
		{"empty token", "", hash},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if ValidateRefreshTokenHash(tc.incoming, tc.dbHash) {
				t.Error("expected validation to fail")
			}
		})
	}
}
