package identity

import (
	"context"
	"fmt"
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
