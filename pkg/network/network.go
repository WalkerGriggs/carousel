package network

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/channel"
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
	Name     string             `json:"name"`
	URI      uri.URI            `json:"uri"`
	Ident    identity.Identity  `json:"ident"`
	Buffer   chan *irc.Message  `json:",omitempty"`
	Channels []*channel.Channel `json:",omitempty"`

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

// Listen attempts to connect and listen if the network isn't already connectedif network isn't already connected and  listens to ensures the network is connected and listening. If the
// network is already connected, this function is a no-op.
//
// If the Network's Connection is nil when Wide is called, Wide will attempt to
// connected the Network. If the connection fails, Wide logs an error and exits.
func (n *Network) Listen() {
	if n.Connection == nil {
		n.LogEntry().Debug("Establishing connection")
		n.Buffer = make(chan *irc.Message)

		err := n.connect()
		if err != nil {
			n.LogEntry().WithError(err).Error("Failed to connect to network.")
			return
		}

		go n.listen()
	} else {
		n.localReply()
	}
}

// listen reads, parses, and forwards all mesages sent from the network to the
// client. in it's current state, this blocking function should exit if the
// network encounters an error when receiving messages.
func (n *Network) listen() {
	n.LogEntry().Debug("Listening to network.")
	for {
		msg, err := n.Receive()
		if err != nil {
			n.LogEntry().WithError(err).Error("Failed to receive message.")
			return
		}

		hook := CommandTable[msg.Command]
		if hook != nil {
			send, err := hook(n, msg)
			if err != nil {
				n.LogEntry().WithError(err).Error(err.Error())
			}

			if !send { continue }

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

func (n *Network) localReply() {
	for _, msg := range n.ClientReplies {
		n.Buffer <- msg
	}
}

func (n *Network) isJoined(name string) bool {
	for _, channel := range n.Channels {
		if channel.Name == name {
			return true
		}
	}
	return false
}

func (n *Network) getChannel(name string) *channel.Channel {
	for _, channel := range n.Channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}
