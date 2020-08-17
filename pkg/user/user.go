package user

import (
	"github.com/pkg/errors"
	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/crypto/phash"
	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
)

type Options struct {
	Username string
	Password string
	Networks []*network.Network
}

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Username string             `json:"username"`
	Password string             `json:"password"`
	Networks []*network.Network `json:"networks,omitempty"`
	Client   *client.Client     `json:",omitempty"`
}

func New(opts Options) (*User, error) {
	return &User{
		Username: opts.Username,
		Password: opts.Password,
		Networks: opts.Networks,
	}, nil
}

// Authorized compares the given password with the password hash stored in the
// config. The user's password isn't stored in plaintext (for very obvious
// reasons, so we have to hash and salt the supplied password before comparing)
func (u *User) Authorize(ident identity.Identity) error {
	if !phash.HashesMatch(u.Password, ident.Password) {
		return errors.Errorf("Authorization for user %s failed.", ident.Username)
	}
	return nil
}

func (u *User) GetNetwork(name string) (*network.Network, error) {
	for _, network := range u.Networks {
		if network.Name == name {
			return network, nil
		}
	}

	return nil, errors.Errorf("Network %s not found for user %s", name, u.Username)
}

func (u *User) GetDefaultNetwork() (*network.Network, error) {
	if len(u.Networks) == 0 {
		return nil, errors.Errorf("User %s has no networks", u.Username)
	}

	return u.Networks[0], nil
}
