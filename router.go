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

	LocalConn *Connection
	WideConn  *irc.Conn

	LocalRpl []*irc.Message
}

func NewRouter() *Router {
	return &Router{
		Local: make(chan *irc.Message),
		Wide:  make(chan *irc.Message),
	}
}

// Read reads, santizies, and forwards all messages sent from the User to the
// network. In its current state, this blocking process should only exit if the
// reader throws an error.
//
// Network <-> (Bouncer <- Client)
func (r *Router) LocalRead() {
	reader := bufio.NewReader(r.LocalConn.Conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		msg = strings.TrimSpace(string(msg))

		if parsed_msg := irc.ParseMessage(msg); parsed_msg != nil {
			switch parsed_msg.Command {

			case "QUIT":
				r.LocalConn = nil
				return

			default:
				r.Wide <- parsed_msg
			}
		}
	}
}

// Read decodes each message sent from the Network and forwards it to the User.
// In its current state, this blocking process should only exit if the decoder
// throws an error.
//
// (Network -> Bouncer) <-> Client
func (r *Router) WideRead() {
	for {
		msg, err := r.WideConn.Decode()
		if err != nil {
			log.Fatal(err)
		}

		switch msg.Command {

		case "PING":
			Pong(r.WideConn, msg)

		case "001", "002", "003", "004", "005":
			r.LocalRpl = append(r.LocalRpl, msg)
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
func (r *Router) WideWrite() {
	for msg := range r.Wide {
		if err := r.WideConn.Encode(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func (r *Router) LocalWrite() {
	for msg := range r.Local {

		// Temp hack until we implement better signal handler
		if r.LocalConn == nil {
			return
		}

		if _, err := r.LocalConn.Conn.Write([]byte(msg.String() + "\n")); err != nil {
			log.Fatal(err)
		}
	}
}

func (r Router) LocalReply() {
	r.localBatchSend(r.LocalRpl)
}

func (r Router) localBatchSend(messages []*irc.Message) {
	for _, msg := range messages {
		if _, err := r.LocalConn.Conn.Write([]byte(msg.String() + "\n")); err != nil {
			log.Fatal(err)
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
