package carousel

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"
)

type NetworkConfig struct {
	Name  string
	URI   string
	Ident *Identity
}

// Network represents an IRC network. Each network has a URI, and, because Users
// own the Network object, each Network stores the User's Identity as well.
type Network struct {
	Config   *NetworkConfig
	Name     string     `json:"name"`
	URI      string     `json:"uri"`
	Ident    *Identity  `json:"ident,omitempty"`
	Channels []*Channel `json:",omitempty"`

	Buffer        chan *irc.Message `json:"-,omitempty"`
	Connection    *irc.Conn         `json:"-,omitempty"`
	ClientReplies []*irc.Message    `json:"-,omitempty"`
}

// New takes in a Network Config and returns a new Network object. In this case,
// all configs are manadatory, but, in its current state, New doesn't throw any
// errors.
func NewNetwork(config *NetworkConfig) (*Network, error) {
	return &Network{
		Config: config,
		Name:   config.Name,
		URI:    config.URI,
		Ident:  config.Ident,
		Buffer: make(chan *irc.Message),
	}, nil
}

// Listen attempts to connect and listen if the Network isn't already connected.
// If the Network's Connection is nil, Listen reads parses, and forwards all
// messages from the Network to the Client. In it's current state, this blocking
// function should exit if the Network encounters an error when receiving
// messages.
func (n *Network) Listen() error {
	if n.Connection != nil {
		n.localReply()
		return nil
	}

	n.Buffer = make(chan *irc.Message)
	if err := n.connect(); err != nil {
		n.LogEntry().WithError(err).Error("Failed to connect to network")
		return err
	}

	for {
		msg, err := n.Receive()
		if err != nil {
			n.LogEntry().WithError(err).Error("Network failed to receive a message")
			return err
		}

		send, err := n.MaybeRun(msg)
		if err != nil {
			n.LogEntry().WithError(err).Error(err)
			return err
		}

		if send {
			n.Buffer <- msg
		}
	}
}

// connect dials the network and identifies. If the dial throws an error,
// connect short circuits -- handle this accordingly.
func (n *Network) connect() error {
	conn, err := irc.Dial(n.URI)
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

// localReply relays welcome messages to the client.
func (n *Network) localReply() error {
	for _, msg := range n.ClientReplies {
		n.Buffer <- msg
	}
	return nil
}

// isJoined checks checks if the network has already joined the give channel.
func (n *Network) isJoined(name string) bool {
	for _, channel := range n.Channels {
		if channel.Name == name {
			return true
		}
	}
	return false
}

// getChannel searches the networks channels and retrieves the channel matching
// the given channel name.
func (n *Network) getChannel(name string) *Channel {
	for _, channel := range n.Channels {
		if channel.Name == name {
			return channel
		}
	}
	return nil
}

func (n *Network) Send(msg *irc.Message) error {
	return n.Connection.Encode(msg)
}

func (n *Network) Receive() (*irc.Message, error) {
	return n.Connection.Decode()
}

func (n *Network) BatchSend(messages []*irc.Message) error {
	for _, msg := range messages {
		if err := n.Send(msg); err != nil {
			return err
		}
	}
	return nil
}

func (n *Network) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"Network": n.Name,
		"Host":    n.URI,
		"User":    n.Ident.Username,
	})
}
