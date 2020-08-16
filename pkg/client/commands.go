package client

import (
	"gopkg.in/sorcix/irc.v2"
)

// CommandHook represents a function which should be run in response to an IRC
// message. A CommandHook returns a boolean, and an error. The boolean is true
// if the message should be intercepted or passed along.
type CommandHook func(n *Client, msg *irc.Message) (bool, error)

// CommandTable maps the command string to a CommandHook.
type CommandTable map[string]CommandHook

// MaybeRun wrap a given CommandTable and runs the given message's corresponding
// hook. If the message command isn't a key in the CommandTable, it does nothing
// and returns true.
func (t CommandTable) MaybeRun(n *Client, msg *irc.Message) (bool, error) {
	hook := t[msg.Command]
	if hook == nil {
		return true, nil
	}

	return hook(n, msg)
}

// CommandTable maps IRC to commands to hooks. Whenever the client receives a
// message, it looks up the command in the CommandTable and runs the corresponding
// function.
var ClientCommandTable = CommandTable{
	"USER": (*Client).user,
	"NICK": (*Client).nick,
	"PASS": (*Client).pass,
	"QUIT": (*Client).quit,
}

// user pulls identity parameters out of the message and stores them in the
// client. This ident is used to authenticate the client as a user, and should
// not be passed to the network.
func (c *Client) user(msg *irc.Message) (bool, error) {
	c.Ident.Username = msg.Params[0]
	c.Ident.Realname = msg.Params[3]
	return false, nil
}

// nick pulls identity parameters out of the message and stores them in the
// client. This ident is used to authenticate the client as a user, and should
// not be passed to the network.
func (c *Client) nick(msg *irc.Message) (bool, error) {
	c.Ident.Nickname = msg.Params[0]
	return false, nil
}

// pass pulls identity parameters out of the message and stores them in the
// client. This ident is used to authenticate the client as a user, and should
// not be passed to the network.
func (c *Client) pass(msg *irc.Message) (bool, error) {
	c.Ident.Password = msg.Params[0]
	return false, nil
}

// quit disconnects the client and returns an ErrDisconnected error which will
// bubble up to the rungroup and cause all client-space routines (heartbeat,
// router, etc) to return as well.
func (c *Client) quit(msg *irc.Message) (bool, error) {
	c.Disconnect()
	return false, ErrDisconnected
}
