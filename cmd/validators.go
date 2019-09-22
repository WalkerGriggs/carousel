package cmd

import (
	"errors"
	"strconv"
)

// Valid password returns and error if password requirements haven't been met.
// TODO more formal password validation
func validatePassword(input string) error {
	if len(input) < 10 || len(input) > 64 {
		return errors.New("Password must be between 10 and 64 characters long")
	}
	return nil
}

// validatePort return an error if the input string doesn't represent an
// integer, or isn't within the valid port range (1 - 65535)
func validatePort(input string) error {
	i, err := strconv.Atoi(input)
	if err != nil || i < 1 || i > 65535 {
		return errors.New("Invalid port number! The port should be an integer between 1 and 65535.")
	}

	return nil
}

// validateNotEmpty return an error if the input string is empty.
func validateNotEmpty(input string) error {
	if input == "" {
		return errors.New("Empty!")
	}
	return nil
}
