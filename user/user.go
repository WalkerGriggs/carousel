package user

import (
	"github.com/walkergriggs/carousel/network"
	"github.com/walkergriggs/carousel/router"
	"github.com/walkergriggs/carousel/crypto/phash"
)

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Network  *network.Network `json:"network,omitempty"`
	Router   *router.Router   `json:",omitempty"`
}

// Authorized compares the given password with the password hash stored in the
// config. The user's password isn't stored in plaintext (for very obvious
// reasons, so we have to hash and salt the supplied password before comparing)
func (u User) Authorized(ident network.Identity) bool {
	return phash.HashesMatch(u.Password, ident.Password)
}
