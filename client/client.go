package client

import (
	"bufio"
	"gopkg.in/sorcix/irc.v2"
	"net"

	"github.com/walkergriggs/carousel/network"
)

// Connection represents the connection between the User and the Server. Each
// connection maintains a network connection, accepted in the Server's main
// event loop.
type Client struct {
	Connection net.Conn          `json:",omitempty"`
	Buffer     chan *irc.Message `json:",omitempty"`
	Ident      *network.Identity `json:",omitempty"`

	encoder *irc.Encoder  `json:",omitempty"`
	decoder *irc.Decoder  `json:",omitempty"`
	reader  *bufio.Reader `json:",omitempty"`
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Connection: conn,
		Buffer:     make(chan *irc.Message),
		Ident:      new(network.Identity),

		encoder: irc.NewEncoder(conn),
		decoder: irc.NewDecoder(conn),
		reader:  bufio.NewReader(conn),
	}
}

func (c *Client) Local() {
	for {
		msg, err := c.Receive()
		if err != nil {
			c.LogEntry().WithError(err).Error("Unable to receive message.")
			return
		}

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

func (c *Client) Ping(nickname string) {
	c.Send(&irc.Message{
		Command: "PING",
		Params:  []string{nickname},
	})
}
