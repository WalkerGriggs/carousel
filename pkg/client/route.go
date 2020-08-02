package client

import (
	"github.com/walkergriggs/carousel/pkg/network"
)

// Route passes messages from the given Network buffer to the Client buffer, and
// visa versa. Route also calls heartbeat to periodically ping the Client's
// Connection. If the Client doesn' respond to the Ping or encounters an error
// when sending to either the Client or Network, Route returns.
func (c *Client) Route(n *network.Network) {
	for {
		select {
		case <-c.disconnect:
			return

		case msg := <-c.Buffer:
			if err := n.Send(msg); err != nil {
				n.LogEntry().WithError(err).Error("Failed to send to network.")
			}

		case msg := <-n.Buffer:
			if err := c.Send(msg); err != nil {
				c.LogEntry().WithError(err).Error("Failed to send to client.")
			}
		}
	}
}
