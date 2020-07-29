package user

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/pkg/identity"
	"github.com/walkergriggs/carousel/pkg/network"
)

// Route passes messages from the given Network buffer to the Client buffer, and
// visa versa. Route also calls heartbeat to periodically ping the Client's
// Connection. If the Client doesn' respond to the Ping or encounters an error
// when sending to either the Client or Network, Route returns.
func (u *User) Route(n *network.Network) {
	timeout := make(chan bool, 1)
	msgs := make(chan *irc.Message)
	go u.heartbeat(n.Ident, msgs, timeout)

	for {
		select {
		case <-timeout:
			u.Client.Connection.Close()
			return

		case msg := <-u.Client.Buffer:
			if msg.Command == "PONG" {
				msgs <- msg
			}

			if err := n.Send(msg); err != nil {
				n.LogEntry().WithError(err).Error("Failed to send to network.")
			}

		case msg := <-n.Buffer:
			if err := u.Client.Send(msg); err != nil {
				u.Client.LogEntry().WithError(err).Error("Failed to send to client.")
			}
		}
	}
}

// heartbeat sends a Ping message ot the Client and waits for a response. If
// it doesn't hear back within a few seconds, heartbeat will send an message
// over the exit channel.
func (u *User) heartbeat(ident identity.Identity, msgs <-chan *irc.Message, exit chan<- bool) {
	// Ping the Client every 30 seconds
	for range time.Tick(30 * time.Second) {
		timeout := make(chan bool)

		u.Client.Ping(ident.Nickname)

		// Send a timeout message after 5 seconds.
		go func(timeout chan<- bool) {
			time.Sleep(5 * time.Second)
			timeout <- true
		}(timeout)

		// If the select receives the timeout message before the Pong message, log
		// an error and return. Otherwise, loop again.
		select {
		case <-msgs:
			continue

		case <-timeout:
			log.WithFields(log.Fields{
				"Nickname": ident.Nickname,
			}).Warn("Failed to receive PONG. Disconnecting client")

			exit <- true
			return
		}
	}
}
