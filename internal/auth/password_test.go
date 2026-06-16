package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if hash == "secret123" {
		t.Error("hash must not equal the original password")
	}
}

func TestHashPassword_UniqueSalts(t *testing.T) {
	h1, _ := HashPassword("secret123")
	h2, _ := HashPassword("secret123")
	if h1 == h2 {
		t.Error("expected different hashes for the same input — bcrypt must generate a random salt each time")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("unexpected error hashing: %v", err)
	}
	if !CheckPasswordHash("secret123", hash) {
		t.Error("expected correct password to match its hash")
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	hash, _ := HashPassword("secret123")
	for _, tc := range []struct {
		name     string
		password string
	}{
		{"wrong password", "wrong"},
		{"empty password", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if CheckPasswordHash(tc.password, hash) {
				t.Errorf("expected %q to not match hash", tc.password)
			}
		})
	}
}
