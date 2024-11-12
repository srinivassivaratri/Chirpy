package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt during user creation/registration
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPasswordHash checks if a password matches its hash during login/authentication
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
