package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password. The bcrypt algorithm automatically generates a salt and includes it in the resulting hash string.
func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost (equal to 10) is the optimal balance between speed and security.
	// The algorithm will generate a random salt and embed it directly into the hash string.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain password with a hashed password from the database. The bcrypt.CompareHashAndPassword function extracts the salt from the hash and uses it to hash the provided password before comparing.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
