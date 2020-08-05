package router

import (
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/network"
)

type Router struct {
	Client      *client.Client
	Network     *network.Network
	HostAddress string
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

func (r *Router) AttachClient() {
	r.Client.LogEntry().Debug("Attaching to existing channels")

	prefix := &irc.Prefix{
		Name: r.Client.Ident.Nickname,
		User: r.Client.Ident.Username,
		Host: r.HostAddress,
	}

	for _, channel := range r.Network.Channels {
		r.Client.Send(&irc.Message{
			Prefix:  prefix,
			Command: "JOIN",
			Params:  []string{channel.Name},
		})

		// TODO: this is a bit of a hack. Instead of sending names to the client,
		// ask the network for channel names which will be routed through to the
		// client.
		r.Network.Send(&irc.Message{
			Command: "NAMES",
			Params:  []string{channel.Name},
		})
	}
}

func (r *Router) DetachClient() {
	r.Client.LogEntry().Debug("Detaching from channels")

	for _, channel := range r.Network.Channels {
		r.Client.Send(&irc.Message{
			Command: "PART",
			Params:  []string{channel.Name},
		})
	}
}
