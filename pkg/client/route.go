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


// heartbeat sends a Ping message ot the Client and waits for a response. If
// it doesn't hear back within a few seconds, heartbeat will send an message
// over the exit channel.
// func (c *Client) heartbeat(ident identity.Identity, msgs <-chan *irc.Message, exit chan<- bool) {
// 	// Ping the Client every 30 seconds
// 	for range time.Tick(30 * time.Second) {
// 		timeout := make(chan bool)

// 		c.Ping(ident.Nickname)

// 		// Send a timeout message after 5 seconds.
// 		go func(timeout chan<- bool) {
// 			time.Sleep(5 * time.Second)
// 			timeout <- true
// 		}(timeout)

// 		// If the select receives the timeout message before the Pong message, log
// 		// an error and return. Otherwise, loop again.
// 		select {
// 		case <-msgs:
// 			continue

// 		case <-timeout:
// 			log.WithFields(log.Fields{
// 				"Nickname": ident.Nickname,
// 			}).Warn("Failed to receive PONG. Disconnecting client")

// 			exit <- true
// 			return
// 		}
// 	}
// }
