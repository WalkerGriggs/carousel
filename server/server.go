package server

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/client"
	"github.com/walkergriggs/carousel/router"
	"github.com/walkergriggs/carousel/uri"
	"github.com/walkergriggs/carousel/user"
)

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	URI             uri.URI      `json:"uri"`
	Users           []*user.User `json:"users"`
	SSLEnabled      bool         `json:"sslEnabled"`
	CertificatePath string       `json:"certificatePath"`
	Listener        net.Listener `json:",omitempty"`
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

	s.loadUsers()

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
	c := client.NewClient(conn)
	c.LogEntry().Debug("Accepting client connection.")

	go c.Local()

	// authorize is a blocking function, and will not return until the user has
	// been authroized or (todo) timeout reached
	u, err := s.authorize(c)
	if err != nil {
		c.LogEntry().WithError(err).Error("Failed to authorize client.")
		return
	}

	u.Router.Client = c

	if u.Router.Network.Connection == nil {
		go u.Router.Network.Wide()
	}

	go u.Router.Route()
	u.Router.LocalReply()
}

func (s Server) authorize(c *client.Client) (*user.User, error) {
	for {
		if c.Ident.Username == "" {
			continue
		}

		u, err := s.authorizeClient(c)
		if err != nil {
			c.LogEntry().WithError(err).Error("Failed to authenticate with user %s. Retrying.\n", c.Ident.username)
		}

		if u != nil {
			c.LogEntry().Infof("Client authenticated with user %s.\n", u.Username)
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

// LoadUsers performs verious pre-flight checks and operations on all Users
// before the Server begins accepting connections.
func (s Server) loadUsers() {
	for _, u := range s.Users {
		if u.Router == nil {
			u.Router = router.NewRouter(nil, u.Network)
			u.Router.Network.Buffer = make(chan *irc.Message)
		}
	}
}

// getUser searches the server's users and retrieves the user matching the given
// username. This function is only a helper until a better User storage solution
// is implemented.
func (s Server) GetUser(username string) *user.User {
	return GetUser(s.Users, username)
}

func GetUser(users []*user.User, username string) *user.User {
	for _, user := range users {
		if username == user.Username {
			return user
		}
	}
	return nil
}
