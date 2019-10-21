package network

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/uri"
)

// Network represents an IRC network. Each network has a URI, and, because Users
// own the Network object, each Network stores the User's Identity as well.
type Network struct {
	Connection *irc.Conn `json:",omitempty"`
	Name       string    `json:"name"`
	URI        uri.URI   `json:"uri"`
	Ident      Identity  `json:"ident"`

	ClientReplies []*irc.Message    `json:",omitempty"`
	Buffer        chan *irc.Message `json:",omitempty"`
}

// Identity represnts the necessary information to authenticate with a Network.
// See RFC 2812 § 3.1
type Identity struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
}

func (n *Network) Connect() error {
	conn, err := irc.Dial(n.URI.String())
	if err != nil {
		return err
	}

	n.Connection = conn
	n.Identify()

	return nil
}

func (n *Network) Wide() {
	err := n.Connect()
	if err != nil {
		n.LogEntry().Error(err)
		return
	}

	for {
		msg, err := n.Receive()
		if err != nil {
			n.LogEntry().Error(err)
			continue
		}

		switch msg.Command {
		case "PING":
			n.Pong(msg)

		case "001", "002", "003", "004", "005":
			n.ClientReplies = append(n.ClientReplies, msg)
		}

		n.Buffer <- msg
	}
}

// Pong responds to the network's Ping with a Pong command.
// See RFC 2812 § 3.7.2
func (n *Network) Pong(msg *irc.Message) {
	n.Send(&irc.Message{
		Command: "PONG",
		Params:  msg.Params,
	})
}

// Identify handles connection registration for each user.
// Again, see RFC 2812 § 3.1
func (n *Network) Identify() {
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

	n.BatchSend(messages)
}
