package carousel

import (
	"fmt"

	"github.com/walkergriggs/carousel/pkg/crypto/phash"
)

type UserConfig struct {
	Username string
	Password string
	Networks []*Network
}

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Config   *UserConfig
	Username string     `json:"username"`
	Password string     `json:"password"`
	Networks []*Network `json:"network,omitempty"`
	Client   *Client    `json:",omitempty"`
}

func NewUser(config *UserConfig) (*User, error) {
	return &User{
		Username: config.Username,
		Password: config.Password,
		Networks: config.Networks,
	}, nil
}

// Authorized compares the given password with the password hash stored in the
// config. The user's password isn't stored in plaintext (for very obvious
// reasons, so we have to hash and salt the supplied password before comparing)
func (u *User) Authorize(ident Identity) error {
	if !phash.HashesMatch(u.Password, ident.Password) {
		return fmt.Errorf("Authorization for user %s failed.", ident.Username)
	}
	return nil
}

// NetworkOrDefault, like Network, returns the network with matching name or
// returns nil if no match if found. If the given string is empty, however, it
// returns the default network (for now, the default is the first in a non-empty
// network list.
func (u *User) NetworkOrDefault(name string) (n *Network) {
	if name != "" {
		n = u.Network(name)
	} else if len(u.Networks) >= 1 {
		n = u.Networks[0]
	}
	return
}

// Network returns the network object with matching name, and returns nil if no
// match is found.
func (u *User) Network(name string) *Network {
	for _, net := range u.Networks {
		if net.Name == name {
			return net
		}
	}
	return nil
}
