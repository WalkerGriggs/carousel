package client

import (
	"bufio"
	"fmt"
	"net"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/identity"
)

type Options struct {
	Connection net.Conn
}

// Client represents the actual Connection between the User and the Server.
// Unlike Network which represents both the Connection and metadata, Client
// is seperate from User so it can run independetly before identifying itself.
type Client struct {
	Connection net.Conn           `json:",omitempty"`
	Buffer     chan *irc.Message  `json:",omitempty"`
	Ident      *identity.Identity `json:",omitempty"`

	Encoder *irc.Encoder  `json:",omitempty"`
	Decoder *irc.Decoder  `json:",omitempty"`
	Reader  *bufio.Reader `json:",omitempty"`
}

// New takes in Client Options and returns a new Client object. In this case,
// the Connection option is actually mandatory, and New will throw an error only
// if that connection is nil.
func New(opts Options) (*Client, error) {
	if opts.Connection == nil {
		return nil, fmt.Errorf("Failed to create user client. No connection provided.")
	}

	return &Client{
		Connection: opts.Connection,
		Buffer:     make(chan *irc.Message),
		Ident:      new(identity.Identity),

		Encoder: irc.NewEncoder(opts.Connection),
		Decoder: irc.NewDecoder(opts.Connection),
		Reader:  bufio.NewReader(opts.Connection),
	}, nil
}

// Local reads, sanitizes, and forwards all messages sent from the User directed
// towards the Network. In its current state, this blocking process should only
// exit if...
//   - the reader throws an error
//   - the Client disconnects
func (c *Client) Local() {
	c.LogEntry().Debug("Listening over client connection.")

	for {
		msg, err := c.Receive()
		if err != nil {
			c.LogEntry().WithError(err).Error("Failed to receive message.")
			return
		}

		// Parse and store USER, NICK, and PASS commands used by the client to
		// authenticate with specified user. Otherwise, pass the message along.x
		switch msg.Command {
		case "USER":
			c.Ident.Username = msg.Params[0]
			c.Ident.Realname = msg.Params[3]

		case "NICK":
			c.Ident.Nickname = msg.Params[0]

		case "PASS":
			c.Ident.Password = msg.Params[0]

		case "QUIT":
			c.LogEntry().Debug("Client disconnected")
			c.Connection.Close()
			return

		default:
			c.Buffer <- msg
		}
	}
}

// Ping sends a simple PING message to the client. See RFC 2812 § 3.7.2 for more
// details.
func (c *Client) Ping(nickname string) {
	c.Send(&irc.Message{
		Command: "PING",
		Params:  []string{nickname},
	})
}
