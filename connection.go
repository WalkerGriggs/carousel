package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"gopkg.in/sorcix/irc.v2"
)

// Connection represents the connection between the User and the Server. Each
// connection maintains a network connection, accepted in the Server's main
// event loop.
type Connection struct {
	Conn    net.Conn
	Encoder *irc.Encoder
	Decoder *irc.Decoder
}

func NewConnection(conn net.Conn) Connection {
	return Connection{
		Conn:    conn,
		Encoder: irc.NewEncoder(conn),
		Decoder: irc.NewDecoder(conn),
	}
}

// decodeIdent receives identification commands from an RFC compliant IRC
// client. It listens for the necessary connection registration commands which
// the User uses to verify against their credentials. These messages are not
// forwarded on to the Network, but consumed by decodeIdent. Connection
// registration is handled by the server using the Identity specified in the
// User's config.
//
// decodeIdent is blocking, and will not return unless all of the required
// commands have been supplied or a timeout has been reached (to be
// implemented).
func (c Connection) decodeIdent(router Router) string {
	messages := make(map[string]*irc.Message)
	required_commands := []string{"USER", "NICK", "PASS"}

	for {
		message, err := c.Decoder.Decode()
		if err != nil {
			log.Fatal(err)
		}

		messages[strings.ToUpper(message.Command)] = message

		if containsAll(messages, required_commands) {
			break
		}
	}

	// TODO: Proper login logic
	return messages["USER"].Params[0]
}

// ReadWrite spawns individual, concurrent processes for the Connection's read
// and write tasks. ReadWrite should only return when both blocking processes
// return.
func (c Connection) ReadWrite(router Router) {
	go c.Read(router)
	go c.Write(router)
}

// Read reads, santizies, and forwards all messages sent from the User to the
// network. In its current state, this blocking process should only exit if the
// reader throws an error.
//
// Network <-> (Bouncer <- Client)
func (c Connection) Read(router Router) {
	reader := bufio.NewReader(c.Conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		msg = strings.TrimSpace(string(msg))

		if irc_message := irc.ParseMessage(msg); irc_message != nil {
			router.Wide <- irc_message
		}
	}
}

// Write writes all messages passed from the IRC network to the User's TCP
// connection. In its current state, this blocking process should only exit if
// the writer throws an error.
//
// Network <-> (Bouncer -> Client)
func (c Connection) Write(router Router) {
	for {
		msg := <-router.Local
		if _, err := c.Conn.Write([]byte(msg.String() + "\n")); err != nil {
			log.Fatal(err)
		}
	}
}

// containsAll checks to see if the message map contains all of the required
// commands.
func containsAll(messages map[string]*irc.Message, required_commands []string) bool {
	for _, command := range required_commands {
		if _, ok := messages[command]; !ok {
			return false
		}
	}

	return true
}
