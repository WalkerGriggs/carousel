package carousel

import (
	"bufio"
	"context"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
)

type ClientConfig struct {
	Connection net.Conn
}

// Client represents the actual Connection between the User and the Server.
// Unlike Network which represents both the Connection and metadata, Client
// is seperate from User so it can run independetly before identifying itself.
type Client struct {
	Connection net.Conn          `json:",omitempty"`
	Buffer     chan *irc.Message `json:",omitempty"`
	Ident      *Identity         `json:",omitempty"`

	Encoder *irc.Encoder  `json:",omitempty"`
	Decoder *irc.Decoder  `json:",omitempty"`
	Reader  *bufio.Reader `json:",omitempty"`
}

// New takes in a ClientConfig and returns a new Client object.
func NewClient(config ClientConfig) (*Client, error) {
	return &Client{
		Connection: config.Connection,
		Buffer:     make(chan *irc.Message),
		Ident:      new(Identity),

		Encoder: irc.NewEncoder(config.Connection),
		Decoder: irc.NewDecoder(config.Connection),
		Reader:  bufio.NewReader(config.Connection),
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

			send, err := c.MaybeRun(msg)
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

func (c *Client) Send(msg *irc.Message) error {
	_, err := c.Connection.Write([]byte(msg.String() + "\n"))
	return err
}

func (c *Client) Receive() (*irc.Message, error) {
	msg, err := c.Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msg = strings.TrimSpace(string(msg))
	return irc.ParseMessage(msg), nil
}

func (c *Client) BatchSend(messages []*irc.Message) error {
	for _, msg := range messages {
		if err := c.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) LogEntry() *log.Entry {
	return log.WithFields(c.LogFields())
}

func (c *Client) LogFields() log.Fields {
	return log.Fields{
		"RemoteAddress": c.Connection.RemoteAddr().String,
	}
}
