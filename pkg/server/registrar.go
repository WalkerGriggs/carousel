package server

import (
	"strings"
	"time"

	"gopkg.in/sorcix/irc.v2"

	"github.com/pkg/errors"
	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/network"
	"github.com/walkergriggs/carousel/pkg/user"
)

func (s Server) blockingAuthorizeClient(c *client.Client) error {
	if err := c.Ident.Wait(30 * time.Second); err != nil {
		return err
	}

	return s.authorizeClient(c)
}

// authorizeClient uses identity information (username and password) provided by
// the connected client to authenticate some user. authorizeClient will return
// an error if provided user does not exist or if the password provided does not
// match the file on record. It returns the user if the credentials are correct.
func (s Server) authorizeClient(c *client.Client) error {
	u, err := s.clientUser(c)
	if err != nil {
		return err
	}

	if err := u.Authorize(*c.Ident); err != nil {
		c.Send(&irc.Message{
			Command: irc.ERR_PASSWDMISMATCH,
			Params:  []string{"irc.carousel.in", c.Ident.Nickname, "Password incorrect"},
		})
		return err
	}

	c.LogEntry().Infof("Client authenticated with user %s.\n", u.Username)
	return nil
}

func (s Server) clientUser(c *client.Client) (*user.User, error) {
	name, err := parseUsername(c.Ident.Username)
	if err != err {
		return nil, err
	}

	return s.GetUser(name)
}

func (s Server) clientNetwork(c *client.Client) (*network.Network, error) {
	name, err := parseNetworkName(c.Ident.Username)
	if err != nil {
		return nil, err
	}

	u, err := s.clientUser(c)
	if err != nil {
		return nil, err
	}

	if name == "" {
		return u.GetDefaultNetwork()
	}

	return u.GetNetwork(name)
}

func parseUsername(s string) (string, error) {
	split := strings.Split(s, "/")
	if len(split) > 2 {
		return "", errors.Errorf("Identity not formatted properly. Usage: <user>/<network>")
	}

	return split[0], nil
}

func parseNetworkName(s string) (string, error) {
	split := strings.Split(s, "/")
	switch len(split) {
	case 2:
		return split[1], nil

	case 1:
		return "", nil

	default:
		return "", errors.Errorf("Identity not formatted properly. Usage: <user>/<network>")
	}
}

// getUser searches the server's users and retrieves the user matching the given
// username. It returns an error if the user does not exist. This function is
// only a helper until a better User storage solution is implemented.
func (s Server) GetUser(username string) (*user.User, error) {
	for _, user := range s.Users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, errors.Errorf("User %s not found", username)
}
