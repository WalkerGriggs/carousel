package network

import (
	"fmt"

	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/channel"
)

// CommandHook represents a function which should be run in response to an IRC
// message. A CommandHook returns a boolean, and an error. The boolean is true
// if the message should be intercepted or passed along.
type CommandHook func(n *Network, msg *irc.Message) (bool, error)

// CommandTable maps the command string to a CommandHook.
type CommandTable map[string]CommandHook

// MaybeRun wrap a given CommandTable and runs the given message's corresponding
// hook. If the message command isn't a key in the CommandTable, it does nothing
// and returns true.
func (t CommandTable) MaybeRun(n *Network, msg *irc.Message) (bool, error) {
	hook := t[msg.Command]
	if hook == nil {
		return true, nil
	}

	return hook(n, msg)
}

// CommandTable maps IRC to commands to hooks. Whenever the network receives a
// message, it looks up the command in the CommandTable and runs the corresponding
// function.
var NetworkCommandTable = CommandTable{
	irc.PING:         (*Network).ping,
	irc.JOIN:         (*Network).join,
	irc.RPL_WELCOME:  (*Network).rpl_welcome,
	irc.RPL_YOURHOST: (*Network).rpl_welcome,
	irc.RPL_CREATED:  (*Network).rpl_welcome,
	irc.RPL_MYINFO:   (*Network).rpl_welcome,
	irc.RPL_BOUNCE:   (*Network).rpl_welcome,
	irc.RPL_NAMREPLY: (*Network).rpl_namreply,
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
