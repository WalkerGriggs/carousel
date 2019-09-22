package carousel

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(s string) (string, error) {
	bytes := []byte(s)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func HashesMatch(hash, s string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s)); err != nil {
		return false
	}

	return true
}
