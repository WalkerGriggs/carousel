package main

import (
	"fmt"
	"log"
	"net"

	"gopkg.in/sorcix/irc.v2"
)

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	URI   URI     `json:"uri"`
	Users []*User `json:"users"`

	Listener net.Listener
}

// URI is the basic information needed to address Networks and Servers. URI is
// not an exhaustive liste of all URI components, and will be extended in future
// implementations.
type URI struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (uri URI) Format() string {
	return fmt.Sprintf("%s:%d", uri.Address, uri.Port)
}

// Serve attaches a tcp listener to the specificed URI, and starts the main
// event loop. Serve blocks for the lifetime of the parent process and should
// only return if the TCP listener closes or errors (even if there are no active
// connections).
func (s Server) Serve() {
	l, err := net.Listen("tcp", s.URI.Format())
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

// Accept establishes a new connection with the accepted TCP client, and spawns
// the concurrent processess responsible to message passing between the IRC
// network and user. Each accepted connection has it's own router and associated
// user, so accept should only return when the user disconnects, or does not
// authenticate.
func (s Server) accept(conn net.Conn) {
	client := NewClient(conn)
	username := client.decodeIdent()

	user := getUser(username, s.Users)

	if user.Router == nil {
		user.Router = &Router{}
	}

	user.Router.Client = &client
	go user.Router.Local()

	if user.Router.IRC == nil {

		wideConn, err := irc.Dial(user.Network.URI.Format())
		if err != nil {
			log.Fatal(err)
		}

		user.Router.IRC = wideConn
		user.Network.Identify(user.Router.IRC)

		go user.Router.Wide()
	}

	user.Router.LocalReply()
}

// getUser searches the server's users and retrieves the user matching the given
// username. This function is only a helper until a better User storage solution
// is implemented.
func getUser(username string, users []*User) *User {
	for _, user := range users {
		if username == user.Username {
			return user
		}
	}

	return nil
}
