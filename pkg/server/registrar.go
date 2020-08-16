package server

import (
	"fmt"
	"time"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/user"
)

func (s Server) blockingAuthorizeClient(c *client.Client) (*user.User, error) {
	if err := c.Ident.Wait(30 * time.Second); err != nil {
		return nil, err
	}

	return s.authorizeClient(c)
}

// authorizeClient uses identity information (username and password) provided by
// the connected client to authenticate some user. authorizeClient will return
// an error if provided user does not exist or if the password provided does not
// match the file on record. It returns the user if the credentials are correct.
func (s Server) authorizeClient(c *client.Client) (*user.User, error) {
	u, err := s.GetUser(c.Ident.Username)
	if err != nil {
		return u, err
	}

	// Ensure the User successfully authorized. If authorization fails, send the
	// client Error 464.
	if err := u.Authorize(*c.Ident); err != nil {
		c.Send(&irc.Message{
			Command: irc.ERR_PASSWDMISMATCH,
			Params:  []string{"irc.carousel.in", c.Ident.Nickname, "Password incorrect"},
		})
		return nil, err
	}

	c.LogEntry().Infof("Client authenticated with user %s.\n", u.Username)
	return u, nil
}

// getUser searches the server's users and retrieves the user matching the given
// username. It returns an error if the user does not exist. This function is
// only a helper until a better User storage solution is implemented.
func (s Server) GetUser(username string) (*user.User, error) {
	for _, user := range s.Users {
		if username == user.Username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("User %s not found", username)
}
