package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hash(s string) string {
	bytes := []byte(s)
	hash, err := bcrypt.GenerateFromPassword(bytes, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

func hashesMatch(hash, s string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(s)); err != nil {
		return false
	}

	return true
}
