package carousel

import (
	"fmt"

	"gopkg.in/sorcix/irc.v2"
)

// NetworkCommandTable maps IRC to commands to hooks. Whenever the client
// receives a message, it looks up the command in the CommandTable and runs the
// corresponding function.
var NetworkCommandTable = map[string]NetworkCommandHook{
	irc.PING:         (*Network).ping,
	irc.JOIN:         (*Network).join,
	irc.RPL_WELCOME:  (*Network).rpl_welcome,
	irc.RPL_YOURHOST: (*Network).rpl_welcome,
	irc.RPL_CREATED:  (*Network).rpl_welcome,
	irc.RPL_MYINFO:   (*Network).rpl_welcome,
	irc.RPL_BOUNCE:   (*Network).rpl_welcome,
	irc.RPL_NAMREPLY: (*Network).rpl_namreply,
}

type NetworkCommandHook func(n *Network, msg *irc.Message) (bool, error)

// CommandTable maps IRC to commands to hooks. Whenever the network receives a
// message, it looks up the command in the CommandTable and runs the corresponding
// function.
func (n *Network) MaybeRun(msg *irc.Message) (bool, error) {
	hook, ok := NetworkCommandTable[msg.Command]
	if !ok {
		return true, nil
	}
	return hook(n, msg)
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
		channel, _ := NewChannel(name)
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
