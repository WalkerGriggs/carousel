package main

import (
	"bufio"
	"log"
	"strings"

	"gopkg.in/sorcix/irc.v2"
)

// Router maintains the two channels over which messages are passed between the
// User and Network. Named after WAN and LAN, the Local channel handles all
// traffic between the User and Router, and the Wide channel handles all traffic
// between the Router and Network
type Router struct {
	Local chan *irc.Message
	Wide  chan *irc.Message

	LocalConn Connection
	WideConn  *irc.Conn
}

func NewRouter() Router {
	return Router{
		Local: make(chan *irc.Message),
		Wide:  make(chan *irc.Message),
	}
}

func (r Router) Route() {
	go r.LocalRead()
	go r.WideRead()
	go r.Write()
}

// Read reads, santizies, and forwards all messages sent from the User to the
// network. In its current state, this blocking process should only exit if the
// reader throws an error.
//
// Network <-> (Bouncer <- Client)
func (r Router) LocalRead() {
	reader := bufio.NewReader(r.LocalConn.Conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		msg = strings.TrimSpace(string(msg))

		if irc_message := irc.ParseMessage(msg); irc_message != nil {
			r.Wide <- irc_message
		}
	}
}

// Read decodes each message sent from the Network and forwards it to the User.
// In its current state, this blocking process should only exit if the decoder
// throws an error.
//
// (Network -> Bouncer) <-> Client
func (r Router) WideRead() {
	for {
		msg, err := r.WideConn.Decode()
		if err != nil {
			log.Fatal(err)
		}

		if msg.Command == "PING" {
			Pong(r.WideConn, msg)
		}

		r.Local <- msg
	}
}

// Write handles both sides (Wide and Local) of the connection. It...
//    - encodes each message sent from the User and sends them off to the
//      network.
//    - writes all messages passed from the IRC network to the User's TCP
//      connection.
//
// In its current state, this blocking process should only exit if the
// encoder throws an error.
func (r Router) Write() {
	for {
		select {
		case msg := <-r.Local:
			if _, err := r.LocalConn.Conn.Write([]byte(msg.String() + "\n")); err != nil {
				log.Fatal(err)
			}

		case msg := <-r.Wide:
			if err := r.WideConn.Encode(msg); err != nil {
				log.Fatal()
			}
		}
	}
}

// Pong responds to the Network's ping with a Pong command.
func Pong(conn *irc.Conn, msg *irc.Message) {
	conn.Encode(&irc.Message{
		Command: "PONG",
		Params:  msg.Params,
	})
}
