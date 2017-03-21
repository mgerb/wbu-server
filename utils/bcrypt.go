package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateHash -
func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	return string(hash), err
}
