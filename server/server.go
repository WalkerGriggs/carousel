package server

import (
	"fmt"
	"log"
	"net"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/client"
	"github.com/walkergriggs/carousel/router"
	"github.com/walkergriggs/carousel/uri"
	"github.com/walkergriggs/carousel/user"
)

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	URI      uri.URI      `json:"uri"`
	Users    []*user.User `json:"users"`
	Listener net.Listener `json:",omitempty"`
}

// Serve attaches a tcp listener to the specificed URI, and starts the main
// event loop. Serve blocks for the lifetime of the parent process and should
// only return if the TCP listener closes or errors (even if there are no active
// connections).
func (s Server) Serve() {
	l, err := net.Listen("tcp", s.URI.String())
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	s.loadUsers()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
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
	// Authorize the user and short circuit if authorization fails
	u, err := s.authorizeConnection(conn)
	if err != nil {
		return
	}

	// Start listening over the local connection.
	go u.Router.Local()

	// Connect and begin listening to the netowkr if not already connected. This
	// should only happen the first time the User connects, the Server should
	// remain connected even after the Client disconnects.
	if u.Router.Network.Connection == nil {
		go u.Router.Wide()
	}

	// Relay necessary connection replies to the user.
	u.Router.LocalReply()
}

// authorizeConnection decodes identity information from client connection and
// authenticates ident against user. If the user exists and authorization is
// successful, authorizeConnection returns the user.  Otherwise,
// authorizeConnection returns an error.
func (s Server) authorizeConnection(conn net.Conn) (*user.User, error) {
	// Get identity information from user. This identity information is used to
	// authenticate with the server -- not the network.
	c := client.NewClient(conn)
	ident := c.DecodeIdent()

	// Ensure the User exists.
	u := s.GetUser(ident.Username)
	if u == nil {
		return nil, fmt.Errorf("User %s not found", ident.Username)
	}

	// Ensure the User successfully authorized. If authorization fails, send the
	// client Error 464.
	if !u.Authorized(ident) {
		c.Send(&irc.Message{
			Command: irc.ERR_PASSWDMISMATCH,
			Params:  []string{"irc.carousel.in", ident.Nickname, "Password incorrect"},
		})

		return nil, fmt.Errorf("Authentication for user %s failed.", ident.Username)
	}

	// If the User exists _and_ they have succesfully authorized, associate the
	// Client, and return the user.
	u.Router.Client = c
	return u, nil
}

// LoadUsers performs verious pre-flight checks and operations on all Users
// before the Server begins accepting connections.
func (s Server) loadUsers() {
	for _, u := range s.Users {
		if u.Router == nil {
			u.Router = router.NewRouter(nil, u.Network)
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
