package server

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/uri"
	"github.com/walkergriggs/carousel/pkg/user"
)

type Options struct {
	URI             uri.URI
	Users           []*user.User
	SSLEnabled      bool
	CertificatePath string
}

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	URI             uri.URI      `json:"uri"`
	Users           []*user.User `json:"users"`
	SSLEnabled      bool         `json:"sslEnabled"`
	CertificatePath string       `json:"certificatePath"`
	Listener        net.Listener `json:",omitempty"`
}

func New(opts Options) (*Server, error) {
	return &Server{
		URI:             opts.URI,
		Users:           opts.Users,
		SSLEnabled:      opts.SSLEnabled,
		CertificatePath: opts.CertificatePath,
	}, nil
}

// Serve attaches a tcp listener to the specificed URI, and starts the main
// event loop. Serve blocks for the lifetime of the parent process and should
// only return if the TCP listener closes or errors (even if there are no active
// connections).
func (s Server) Serve() {
	l, err := s.listener()
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
		}

		go s.accept(conn)
	}
}

// Accept establishes a new connection with the accepted TCP client, and spawns
// the concurrent processess responsible to message passing between the IRC
// network and user. Each accepted connection has it's own router and associated
// user, so accept should only return when the user disconnects, or does not
// authenticate.
func (s Server) accept(conn net.Conn) {
	c, err := client.New(client.Options{
		Connection: conn,
	})
	if err != nil {
		c.LogEntry().Error(err)
		return
	}

	c.Listen()

	// authorize is a blocking function, and will not return until the user has
	// been authroized or (todo) timeout reached
	u, err := s.authorize(c)
	if err != nil {
		c.LogEntry().WithError(err).Error("Failed to authorize client.")
		return
	}

	go c.Route(u.Network)
	u.Network.Listen()
	c.AttachNetwork(u.Network)
}

func (s Server) authorize(c *client.Client) (*user.User, error) {
	for {
		if c.Ident.Username == "" {
			continue
		}

		u, err := s.authorizeClient(c)
		if err != nil {
			c.LogEntry().
				WithError(err).
				Error("Failed to authenticate with user %s. Retrying.\n", c.Ident.Username)
		}

		if u != nil {
			c.LogEntry().Infof("Client authenticated with user %s.\n", u.Username)
			u.Client = c
			return u, nil
		}

		// todo: exponential delay?
		// todo: timeout error?
		time.Sleep(100 * time.Millisecond)
	}
}

// authorizeConnection decodes identity information from client connection and
// authenticates ident against user. If the user exists and authorization is
// successful, authorizeConnection returns the user.  Otherwise,
// authorizeConnection returns an error.
func (s Server) authorizeClient(c *client.Client) (*user.User, error) {
	// Ensure the User exists.
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
