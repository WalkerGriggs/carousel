package user

import (
	"fmt"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/crypto/phash"
	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
)

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Config   *Config
	Username string           `json:"username"`
	Password string           `json:"password"`
	Network  *network.Network `json:"network,omitempty"`
	Client   *client.Client   `json:",omitempty"`
}

func New(config *Config) (*User, error) {
	return &User{
		Username: config.Username,
		Password: config.Password,
		Network:  config.Network,
	}, nil
}

// Authorized compares the given password with the password hash stored in the
// config. The user's password isn't stored in plaintext (for very obvious
// reasons, so we have to hash and salt the supplied password before comparing)
func (u *User) Authorize(ident identity.Identity) error {
	if !phash.HashesMatch(u.Password, ident.Password) {
		return fmt.Errorf("Authorization for user %s failed.", ident.Username)
	}
	return nil
}
