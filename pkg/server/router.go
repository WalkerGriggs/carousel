package server

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/network"
)

type Router struct {
	Server  *Server
	Client  *client.Client
	Network *network.Network
}

// Route passes messages from the given Network buffer to the Client buffer, and
// visa versa. Route also calls heartbeat to periodically ping the Client's
// Connection. If the Client doesn' respond to the Ping or encounters an error
// when sending to either the Client or Network, Route returns.
func (r *Router) Route(done chan bool) {
	r.Client.LogEntry().Debug("Routing messages between client and network")

	for {
		select {
		case <-done:
			r.DetachClient()
			return

		case msg := <-r.Client.Buffer:
			if err := r.Network.Send(msg); err != nil {
				r.Network.LogEntry().WithError(err).Error("Failed to send to network.")
			}

		case msg := <-r.Network.Buffer:
			if err := r.Client.Send(msg); err != nil {
				r.Client.LogEntry().WithError(err).Error("Failed to send to client.")
			}
		}
	}
}

// AttachClient should run after a client first connects and authenticates. It
// joins channels, sets ident, and performs general post-connection tasks.
func (r *Router) AttachClient() {
	r.setIdent()
	r.joinChannels()
}

// DetachClient should run before a client disconnects. Specifically, it parts
// from all channels. The network should stay attached to the channels even after
// the client disconnects, so this messages should be sent directly from the server
// to the client.
func (r *Router) DetachClient() {
	r.Client.LogEntry().Debug("Detaching from channels")

	for _, channel := range r.Network.Channels {
		r.Client.Send(&irc.Message{
			Command: "PART",
			Params:  []string{channel.Name},
		})
	}
}

func (r *Router) setIdent() {
	r.Client.LogEntry().Debug("Setting ident")
	r.Client.Ident = r.Network.Ident
}

// joinChannel sends a join reply, followed by the names list directly to the
// the client. This should only be used directly after the client connects.
func (r *Router) joinChannels() {
	r.Client.LogEntry().Debug("Attaching to existing channels")

	prefix := &irc.Prefix{
		Name: r.Client.Ident.Nickname,
		User: r.Client.Ident.Username,
		Host: r.Server.URI.Host,
	}

	for _, channel := range r.Network.Channels {
		r.Client.Send(&irc.Message{
			Prefix:  prefix,
			Command: "JOIN",
			Params:  []string{channel.Name},
		})

		params := append([]string{prefix.Name, "=", channel.Name}, channel.Nicks...)
		r.Network.LogEntry().Debug(params)
		r.Client.Send(&irc.Message{
			Command: "353",
			Params:  params,
		})

		r.Client.Send(&irc.Message{
			Command: "366",
			Params:  []string{channel.Name, ":End of NAMES list"},
		})
	}
}
