package network

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/channel"
)

type CommandHook func(n *Network, msg *irc.Message)

var CommandTable = map[string]CommandHook{
	"001":  (*Network).rpl_welcome,
	"002":  (*Network).rpl_welcome,
	"003":  (*Network).rpl_welcome,
	"004":  (*Network).rpl_welcome,
	"005":  (*Network).rpl_welcome,
	"PING": (*Network).pong,
	"JOIN": (*Network).join,
	"353":  (*Network).rpl_namreply,
}

func (n *Network) ping(msg *irc.Message) {
	n.pong(msg)
}

func (n *Network) join(msg *irc.Message) {
	name := msg.Params[0]

	if !n.isJoined(name) {
		channel, _ := channel.New(name)
		n.Channels = append(n.Channels, channel)
	}
}

func (n *Network) rpl_welcome(msg *irc.Message) {
	n.ClientReplies = append(n.ClientReplies, msg)
}

func (n *Network) rpl_namreply(msg *irc.Message) {
	channel := n.getChannel(msg.Params[2])
	if channel == nil {
		n.LogEntry().Error("Received names reply for channel which isn't joined")
		return
	}

	channel.AddNicks(msg.Params[3:])
	n.LogEntry().Debug(channel.Nicks)
}
