package client

import (
	"bufio"
	"fmt"
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

	disconnect chan bool `json:",omitempty"`
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

		disconnect: make(chan bool),
	}, nil
}

func (c *Client) Listen() {
	go c.listen()
	go c.heartbeat()
}

// Local reads, sanitizes, and forwards all messages sent from the User directed
// towards the Network. In its current state, this blocking process should only
// exit if...
//   - the reader throws an error
//   - the Client disconnects
func (c *Client) listen() {
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
		case "USER", "NICK", "PASS":
			c.parseIdent(msg)

		case "QUIT":
			c.Disconnect()
			return

		default:
			c.Buffer <- msg
		}
	}
}

func (c *Client) heartbeat() {
	c.LogEntry().Debug("Starting heartbeat for client connection.")

	for range time.Tick(30 * time.Second) {
		select {
		case <-c.disconnect:
			return
		default:
			c.Ping(c.Ident.Nickname)
		}
	}
}

func (c *Client) Disconnect() {
	c.LogEntry().Debug("Client disconnected")
	c.DetachNetwork()
	c.Connection.Close()
	close(c.disconnect)
}

func (c *Client) AttachNetwork(net *network.Network) {
	c.LogEntry().Debug("Attaching to existing channels")
	c.Network = net
	for _, channel := range net.Channels {
		c.Send(&irc.Message{
			Prefix:  c.MessagePrefix(),
			Command: "JOIN",
			Params:  []string{channel.Name},
		})
	}
}

func (c *Client) DetachNetwork() {
	c.LogEntry().Debug("Attaching to existing channels")
	for _, channel := range c.Network.Channels {
		c.Send(&irc.Message{
			Command: "PART",
			Params:  []string{channel.Name},
		})
	}
}

// parseIdent pulls identity parameters out of irc messages and stores them
// in the client.
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
