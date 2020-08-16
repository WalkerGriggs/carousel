package client

import (
	"bufio"
	"context"
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

// Listen reads, sanitizes, and forwards all messages sent from the Client
// directed towards the Network. In its current state, this blocking process
// should only exit if...
//   - the reader throws an error
//   - the Client disconnects
func (c *Client) Listen(ctx context.Context) error {
	c.LogEntry().Debug("Starting to listen to client connection.")

	for {
		select {
		case <-ctx.Done():
			c.LogEntry().Debug("Stopping to listen to client connection.")
			return ctx.Err()

		default:
			msg, err := c.Receive()
			if err != nil {
				return err
			}

			send, err := ClientCommandTable.MaybeRun(c, msg)
			if err != nil {
				return err
			} else if send {
				c.Buffer <- msg
			}
		}
	}
}

// heartbeat sends a ping message every 30 seonds to the client. It takes a done
// channel as a replacement to context.
// TODO: Close the `done` channel if the client doesn't reply with a PONG.
func (c *Client) Heartbeat(ctx context.Context) error {
	c.LogEntry().Debug("Starting client heartbeat.")

	for range time.Tick(30 * time.Second) {
		select {
		case <-ctx.Done():
			c.LogEntry().Debug("Stopping client heartbeat")
			return ctx.Err()

		default:
			c.Ping(c.Ident.Nickname)
		}
	}

	return nil
}

// Disconnect closes the client connection.
func (c *Client) Disconnect() {
	c.Connection.Close()
}

// Ping sends a simple PING message to the client. See RFC 2812 § 3.7.2 for more
// details.
func (c *Client) Ping(nickname string) {
	c.Send(&irc.Message{
		Command: "PING",
		Params:  []string{nickname},
	})
}
