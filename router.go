package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

type Connection interface {
	Send(msg *irc.Message) error
	Receive() (*irc.Message, error)
}

// Router maintains the two channels over which messages are passed between the
// User and Network. Named after WAN and LAN, the Local channel handles all
// traffic between the User and Router, and the Wide channel handles all traffic
// between the Router and Network
type Router struct {
	Client  *Client
	Network *Network

	ClientReplies []*irc.Message
}

func NewRouter(client *Client, network *Network) *Router {
	return &Router{
		Network: network,
		Client:  client,
	}
}

// Local reads, sanitizes, and forwards all messages sent from the User to the
// network. In its current state, this blocking process should exit if...
//   - the reader throws an error
//   - the encoder throws an error
//   - the client disconnects
func (r *Router) Local() error {
	for {
		//msg, err := reader.ReadString('\n')
		msg, err := r.Client.Receive()
		if err != nil {
			return err
		}

		if msg != nil {
			switch msg.Command {
			case "QUIT":
				r.Client = nil
				return nil

			default:
				if err := r.Network.Send(msg); err != nil {
					return err
				}
			}
		}
	}
}

// Wide reads, parses, and forwards all messages send from the network to the
// user. In it's current state, this blocking process should exit if...
//   - the decoder throws an error
//   - the writer throws an error
func (r *Router) Wide() error {
	for {
		//msg, err := r.IRC.Decode()
		msg, err := r.Network.Receive()
		if err != nil {
			return nil
		}

		switch msg.Command {
		case "PING":
			r.Network.Pong(msg)

		case "001", "002", "003", "004", "005":
			r.ClientReplies = append(r.ClientReplies, msg)
		}

		if err := r.Client.Send(msg); err != nil {
			return err
		}
	}
}

// LocalReply relays the reply commands (WELCOME, YOURHOST, CREATED, MYINFO, and
// BOUNCE) initially sent by the network to the user.
// See RFC 2813 § 5.2.1
func (r Router) LocalReply() {
	for _, msg := range r.ClientReplies {
		if _, err := r.Client.Connection.Write([]byte(msg.String() + "\n")); err != nil {
			log.Fatal(err)
		}
	}
}
