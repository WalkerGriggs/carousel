package network

import (
	"fmt"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/channel"
)

type CommandHook func(n *Network, msg *irc.Message) (bool, error)

// CommandTable maps IRC to commands to hooks. Whenever the network receives a
// message, it looks up the command in the CommandTable and runs the corresponding
// function.
var CommandTable = map[string]CommandHook{
	"001":  (*Network).rpl_welcome,
	"002":  (*Network).rpl_welcome,
	"003":  (*Network).rpl_welcome,
	"004":  (*Network).rpl_welcome,
	"005":  (*Network).rpl_welcome,
	"PING": (*Network).ping,
	"JOIN": (*Network).join,
	"353":  (*Network).rpl_namreply,
}

// ping responds to the network with a pong.
func (n *Network) ping(msg *irc.Message) (bool, error) {
	n.pong(msg)
	return false, nil
}

// join adds the specified Channel to the Network.
func (n *Network) join(msg *irc.Message) (bool, error) {
	name := msg.Params[0]
	if !n.isJoined(name) {
		channel, _ := channel.New(name)
		n.Channels = append(n.Channels, channel)
	}
	return true, nil
}

// rpl_welcome records welcome messages to be relayed to the Client on subsequent
// connections.
func (n *Network) rpl_welcome(msg *irc.Message) (bool, error) {
	n.ClientReplies = append(n.ClientReplies, msg)
	return true, nil
}

// rpl_namreply adds specific nicks to a Channel.
func (n *Network) rpl_namreply(msg *irc.Message) (bool, error) {
	channel := n.getChannel(msg.Params[2])
	if channel == nil {
		err := fmt.Errorf("Received name reply for channel which isn't joined.")
		return false, err
	}

	channel.AddNicks(msg.Params[3:])
	return true, nil
}
