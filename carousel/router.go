package carousel

import (
	"context"
	"net/url"

	"gopkg.in/sorcix/irc.v2"
)

type Router struct {
	ServerURI string
	Client    *Client
	Network   *Network
}

// Route passes messages from the given Network buffer to the Client buffer, and
// visa versa. Route also calls heartbeat to periodically ping the Client's
// Connection. If the Client doesn' respond to the Ping or encounters an error
// when sending to either the Client or Network, Route returns.
func (r *Router) Route(ctx context.Context) error {
	r.Client.LogEntry().Debug("Starting to route messages.")

	for {
		select {
		case <-ctx.Done():
			r.Client.LogEntry().Debug("Stopping to route messages.")
			return ctx.Err()

		case msg := <-r.Client.Buffer:
			if err := r.Network.Send(msg); err != nil {
				r.Network.LogEntry().WithError(err).Error("Router failed to send to network.")
			}

		case msg := <-r.Network.Buffer:
			if err := r.Client.Send(msg); err != nil {
				r.Client.LogEntry().WithError(err).Error("Router failed to send to client.")
			}
		}
	}
}

// attachClient should run after a client first connects and authenticates. It
// joins channels, sets ident, and performs general post-connection tasks.
func (r *Router) attachClient() {
	r.setIdent()

	if err := r.joinChannels(); err != nil {
		r.Client.LogEntry().WithError(err).Error(err)
	}
}

// detachClient should run before a client disconnects. Specifically, it parts
// from all channels. The network should stay attached to the channels even after
// the client disconnects, so this messages should be sent directly from the server
// to the client.
func (r *Router) detachClient() error {
	var messages []*irc.Message

	for _, channel := range r.Network.Channels {
		messages = append(messages, &irc.Message{
			Command: "PART",
			Params:  []string{channel.Name},
		})
	}

	return r.Client.BatchSend(messages)
}

func (r *Router) setIdent() {
	r.Client.Ident = r.Network.Ident
}

// joinChannel sends a join reply, followed by the names list directly to the
// the client. This should only be used directly after the client connects.
func (r *Router) joinChannels() error {
	u, err := url.Parse("//" + r.ServerURI)
	if err != nil {
		return err
	}

	prefix := &irc.Prefix{
		Name: r.Client.Ident.Nickname,
		User: r.Client.Ident.Username,
		Host: u.Host,
	}

	var messages []*irc.Message

	for _, channel := range r.Network.Channels {
		messages = append(messages, &irc.Message{
			Prefix:  prefix,
			Command: "JOIN",
			Params:  []string{channel.Name},
		})

		messages = append(messages, &irc.Message{
			Command: "353",
			Params:  append([]string{prefix.Name, "=", channel.Name}, channel.Nicks...),
		})

		messages = append(messages, &irc.Message{
			Command: "366",
			Params:  []string{channel.Name, ":End of NAMES list"},
		})
	}

	return r.Client.BatchSend(messages)
}
