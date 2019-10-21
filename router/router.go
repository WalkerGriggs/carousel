package router

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/sorcix/irc.v2"

	"github.com/walkergriggs/carousel/client"
	"github.com/walkergriggs/carousel/network"
)

type Connection interface {
	Send(msg *irc.Message) error
	Receive() (*irc.Message, error)
	LogWithFields() *log.Entry
}

// Router maintains the two channels over which messages are passed between the
// User and Network. Named after WAN and LAN, the Local channel handles all
// traffic between the User and Router, and the Wide channel handles all traffic
// between the Router and Network
type Router struct {
	Client  *client.Client   `json:",omitempty"`
	Network *network.Network `json:",omitempty"`
}

func NewRouter(client *client.Client, network *network.Network) *Router {
	return &Router{
		Network: network,
		Client:  client,
	}
}

func (r Router) Route() {
	timeout := make(chan bool, 1)
	msgs := make(chan *irc.Message)
	go r.healthcheck(msgs, timeout)

	for {
		select {
		case <-timeout:
			r.Client.Connection.Close()
			return

		case msg := <-r.Client.Buffer:
			if msg.Command == "PONG" {
				msgs <- msg
			}

			if err := r.Network.Send(msg); err != nil {
				r.Network.LogEntry().WithErr(err).Error("Failed to send to network.")
			}

		case msg := <-r.Network.Buffer:
			if err := r.Client.Send(msg); err != nil {
				r.Client.LogEntry().WithErr(err).Error("Failed to send to client.")
			}
		}
	}
}

func (r Router) healthcheck(msgs <-chan *irc.Message, exit chan<- bool) {
	for range time.Tick(30 * time.Second) {
		timeout := make(chan bool)

		r.Client.Ping(r.Network.Ident.Nickname)

		go func(timeout chan<- bool) {
			time.Sleep(5 * time.Second)
			timeout <- true
		}(timeout)

		select {
		case <-msgs:
			continue

		case <-timeout:
			log.WithFields(log.Fields{
				"Nickname": r.Network.Ident.Nickname,
			}).Warn("Failed to receive PONG. Disconnecting client")

			exit <- true
			return
		}
	}
}

// LocalReply relays the reply commands (WELCOME, YOURHOST, CREATED, MYINFO, and
// BOUNCE) initially sent by the network to the user.
// See RFC 2813 § 5.2.1
func (r Router) LocalReply() {
	if err := r.Client.BatchSend(r.Network.ClientReplies); err != nil {
		r.Client.LogEntry().WithError(err).Error("Failed to send to client.")
	}
}
