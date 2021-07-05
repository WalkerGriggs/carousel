package identity

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// Identity represnts the necessary information to authenticate with a Network.
// See RFC 2812 § 3.1
type Identity struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
}

func (i *Identity) CanAuthenticate() bool {
	return i.Username != "" && i.Password != ""

}

// HasNetwork returns true if the username is username/network pair.
func (i *Identity) HasNetwork() bool {
	return strings.ContainsRune(i.Username, '/')
}

// ParsedNetwork returns the string following a '/' in the username. This method
// is naive and assumes that a username will only contain a single slash, and
// the substring following the delimeter is the desired network.
func (i *Identity) ParsedNetwork() string {
	if i.HasNetwork() {
		return strings.Split(i.Username, "/")[1]
	}
	return ""
}

// ParsedUsername returns the string preceding a '/' in the username. This
// method is naive and assumes the desired username is the first substring after
// splitting by a forward slash delimeter. This method should be used instead of
// accessing the field directly.
func (i *Identity) ParsedUsername() string {
	if i.HasNetwork() {
		return strings.Split(i.Username, "/")[0]
	}
	return i.Username
}

// Wait returns if the identity is populated (has a username and password). If
// the ident cannot be authenticated by some given timeout duration, it returns
// an error.
func (i *Identity) Wait(duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	for {
		if i.CanAuthenticate() {
			return nil
		} else if ctx.Err() != nil {
			return fmt.Errorf("Identity not sent after %s", duration.String())
		}
	}
}
