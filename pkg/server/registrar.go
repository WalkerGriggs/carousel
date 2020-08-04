package server

import (
	"fmt"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/user"
)

// authorizeConnection decodes identity information from client connection and
// authenticates ident against user. If the user exists and authorization is
// successful, authorizeConnection returns the user.  Otherwise,
// authorizeConnection returns an error.
func (s Server) authorizeClient(c *client.Client) (*user.User, error) {
	u := s.GetUser(c.Ident.Username)
	if u == nil {
		return nil, fmt.Errorf("User %s not found", c.Ident.Username)
	}

	// Ensure the User successfully authorized. If authorization fails, send the
	// client Error 464.
	if !u.Authorized(*c.Ident) {
		c.Send(&irc.Message{
			Command: irc.ERR_PASSWDMISMATCH,
			Params:  []string{"irc.carousel.in", c.Ident.Nickname, "Password incorrect"},
		})
		return nil, fmt.Errorf("Authentication for user %s failed.", c.Ident.Username)
	}

	c.LogEntry().Infof("Client authenticated with user %s.\n", u.Username)
	return u, nil
}

// getUser searches the server's users and retrieves the user matching the given
// username. This function is only a helper until a better User storage solution
// is implemented.
func (s Server) GetUser(username string) *user.User {
	for _, user := range s.Users {
		if username == user.Username {
			return user
		}
	}
	return nil
}
