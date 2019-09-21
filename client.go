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
type Client struct {
	Connection net.Conn
	Encoder    *irc.Encoder
	Decoder    *irc.Decoder
	Reader     *bufio.Reader
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Connection: conn,
		Reader:     bufio.NewReader(conn),
		Encoder:    irc.NewEncoder(conn),
		Decoder:    irc.NewDecoder(conn),
	}
}

func (c Client) Send(msg *irc.Message) error {
	_, err := c.Connection.Write([]byte(msg.String() + "\n"))
	return err
}

func (c Client) Receive() (*irc.Message, error) {
	msg, err := c.Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msg = strings.TrimSpace(string(msg))

	return irc.ParseMessage(msg), nil
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
func (c Client) decodeIdent() Identity {
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

	return parseIdent(messages)
}

func parseIdent(messages map[string]*irc.Message) Identity {
	return Identity{
		Nickname: messages["NICK"].Params[0],
		Username: messages["USER"].Params[0],
		Realname: messages["USER"].Params[3],
		Password: messages["PASS"].Params[0],
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
