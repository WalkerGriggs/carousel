package user

import (
	"fmt"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/crypto/phash"
	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
)

type Options struct {
	Username string
	Password string
	Network  *network.Network
}

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Network  *network.Network `json:"network,omitempty"`
	Client   *client.Client   `json:",omitempty"`
}

func New(opts Options) (*User, error) {
	return &User{
		Username: opts.Username,
		Password: opts.Password,
		Network:  opts.Network,
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
