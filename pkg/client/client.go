package client

import (
	"bufio"
	"net"
	"time"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
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
	Network    *network.Network   `json:",omitempty"`

	Encoder *irc.Encoder  `json:",omitempty"`
	Decoder *irc.Decoder  `json:",omitempty"`
	Reader  *bufio.Reader `json:",omitempty"`
}

// New takes in Client Options and returns a new Client object.
func New(opts Options) (*Client, error) {
	return &Client{
		Connection: opts.Connection,
		Buffer:     make(chan *irc.Message),
		Ident:      new(identity.Identity),

		Encoder: irc.NewEncoder(opts.Connection),
		Decoder: irc.NewDecoder(opts.Connection),
		Reader:  bufio.NewReader(opts.Connection),
	}, nil
}

func (c *Client) Listen(done chan bool) {
	go c.listen(done)
	go c.heartbeat(done)
}

// Listen reads, sanitizes, and forwards all messages sent from the Client
// directed towards the Network. In its current state, this blocking process
// should only exit if...
//   - the reader throws an error
//   - the Client disconnects
func (c *Client) listen(done chan bool) {
	c.LogEntry().Debug("Listening over client connection.")

	for {
		msg, err := c.Receive()
		if err != nil {
			c.LogEntry().WithError(err).Error("Failed to receive message.")
			return
		}

		if msg.Command == "QUIT" {
			c.Disconnect()
			close(done)
			return
		}

		send, err := ClientCommandTable.MaybeRun(c, msg)
		if err != nil {
			c.LogEntry().WithError(err).Error(err)
			return
		}

		if send {
			c.Buffer <- msg
		}
	}
}

// heartbeat sends a ping message every 30 seonds to the client. It takes a done
// channel as a replacement to context.
// TODO: Close the `done` channel if the client doesn't reply with a PONG.
func (c *Client) heartbeat(done chan bool) {
	c.LogEntry().Debug("Starting heartbeat for client connection.")

	for range time.Tick(30 * time.Second) {
		select {
		case <-done:
			return
		default:
			c.Ping(c.Ident.Nickname)
		}
	}
}

// Disconnect closes the client connection.
func (c *Client) Disconnect() {
	c.LogEntry().Debug("Client disconnected")
	c.Connection.Close()
}

// parseIdent pulls identity parameters out of irc messages and stores them
// in the client. This ident information is used to authenticate the client as a
// user, not to authenticate the client with the network.
func (c *Client) parseIdent(msg *irc.Message) {
	switch msg.Command {
	case "USER":
		c.Ident.Username = msg.Params[0]
		c.Ident.Realname = msg.Params[3]

	case "NICK":
		c.Ident.Nickname = msg.Params[0]

	case "PASS":
		c.Ident.Password = msg.Params[0]
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
