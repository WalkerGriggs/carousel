package server

import (
	"crypto/tls"
	"log"
	"net"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/client"
	"github.com/walkergriggs/carousel/router"
	"github.com/walkergriggs/carousel/crypto/ssl"
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

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.accept(conn)
	}
}

func (s Server) listener() (net.Listener, error) {
	if !s.SSLEnabled {
		return net.Listen("tcp", s.URI.String())
	}

	cert, err := ssl.LoadPem(s.CertificatePath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{*cert}}
	return tls.Listen("tcp", s.URI.String(), config)
}

// Accept establishes a new connection with the accepted TCP client, and spawns
// the concurrent processess responsible to message passing between the IRC
// network and user. Each accepted connection has it's own router and associated
// user, so accept should only return when the user disconnects, or does not
// authenticate.
func (s Server) accept(conn net.Conn) {

	// Get identity information from user. This identity information is used to
	// authenticate with the server -- not the network.
	c := client.NewClient(conn)
	ident := c.DecodeIdent()

	// Get the user the connecting client is authenticating against.
	u := s.GetUser(ident.Username)

	// If the authentication fails, send them err 464 and short circuit
	if !u.Authorized(ident) {
		c.Send(&irc.Message{
			Command: irc.ERR_PASSWDMISMATCH,
			Params:  []string{"irc.carousel.in", ident.Nickname, "Password incorrect"},
		})

		return
	}

	if u.Router == nil {
		u.Router = router.NewRouter(nil, u.Network)
	}

	// Attach the Client connection to the User's Router
	u.Router.Client = c
	go u.Router.Local()

	// Connect to the network if not already connected. This should only happen
	// the first time the User connects, the Server should remain connected even
	// after the Client disconnects.
	if u.Router.Network.Connection == nil {
		err := u.Router.Network.Connect()
		if err != nil {
			log.Fatal(err)
		}

		go u.Router.Wide()
	}

	u.Router.LocalReply()
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
