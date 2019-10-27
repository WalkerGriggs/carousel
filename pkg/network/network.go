package network

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/uri"
)

type Options struct {
	Name  string
	URI   uri.URI
	Ident identity.Identity
}

// Network represents an IRC network. Each network has a URI, and, because Users
// own the Network object, each Network stores the User's Identity as well.
type Network struct {
	Name   string            `json:"name"`
	URI    uri.URI           `json:"uri"`
	Ident  identity.Identity `json:"ident"`
	Buffer chan *irc.Message `json:",omitempty"`

	Connection    *irc.Conn      `json:",omitempty"`
	ClientReplies []*irc.Message `json:",omitempty"`
}

// New takes in Network Options and returns a new Network object. In this case,
// all options are manadatory, but, in its current state, New doesn't throw any
// errors.
func New(opts Options) (*Network, error) {
	return &Network{
		Name:   opts.Name,
		URI:    opts.URI,
		Ident:  opts.Ident,
		Buffer: make(chan *irc.Message),
	}, nil
}

// Wide reads, parses, and forwards all messages send from the network to the
// user. In it's current state, this blocking process should exit if the network
// encounters an error when receiving messages.
//
// If the Network's Connection is nil when Wide is called, Wide will attempt to
// connected the Network. If the connection fails, Wide logs an error and exits.
func (n *Network) Wide() {
	if err := n.connect(); err != nil {
		n.LogEntry().WithError(err).Error("Failed to connect to network.")
		return
	}

	for {
		msg, err := n.Receive()
		if err != nil {
			n.LogEntry().WithError(err).Error("Failed to receive message.")
			return
		}

		switch msg.Command {
		case "PING":
			n.pong(msg)

		case "001", "002", "003", "004", "005":
			n.ClientReplies = append(n.ClientReplies, msg)
		}

		n.Buffer <- msg
	}
}

// connect dials the network and identifies. If the dial throws an error,
// connect short circuits -- handle this accordingly.
func (n *Network) connect() error {
	conn, err := irc.Dial(n.URI.String())
	if err != nil {
		return err
	}

	n.Connection = conn

	n.identify()
	return nil
}

// Identify handles connection registration for each user.
// Again, see RFC 2812 § 3.1
func (n *Network) identify() error {
	var messages []*irc.Message

	if n.Ident.Password != "" {
		messages = append(messages, &irc.Message{
			Command: irc.PASS,
			Params:  []string{n.Ident.Password},
		})
	}

	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params:  []string{n.Ident.Nickname},
	})

	messages = append(messages, &irc.Message{
		Command: irc.USER,
		Params:  []string{n.Ident.Username, "0", "*", n.Ident.Realname},
	})

	return n.BatchSend(messages)
}

// Pong responds to the network's Ping with a Pong command.
// See RFC 2812 § 3.7.2
func (n *Network) pong(msg *irc.Message) {
	n.Send(&irc.Message{
		Command: "PONG",
		Params:  msg.Params,
	})
}
