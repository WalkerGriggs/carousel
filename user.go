package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Network  Network `json:"network"`

	Router Router
}

// ReadWrite connects the user to their Network, and spawns individual,
// concurrent processes for the User's read and write tasks. This function
// should only exit if either Read or Write returns.
func (u User) ReadWrite() {
	conn, err := irc.Dial(u.Network.URI.Format())
	if err != nil {
		log.Fatal(err)
	}

	go u.Network.Identify(conn)
	go u.Read(conn)
	go u.Write(conn)
}

// Read decodes each message sent from the Network and forwards it to the User.
// In its current state, this blocking process should only exit if the decoder
// throws an error.
//
// (Network -> Bouncer) <-> Client
func (u User) Read(conn *irc.Conn) {
	for {
		msg, err := conn.Decode()
		if err != nil {
			log.Fatal(err)
		}

		if msg.Command == "PING" {
			Pong(conn, msg)
		}

		u.Router.Local <- msg
	}
}

// Write encodes each message sent from the User and sends them off to the
// network. In its current state, this blocking process should only exit if the
// encoder throws an error.
//
// (Network <- Bouncer) <-> Client
func (u User) Write(conn *irc.Conn) {
	for msg := range u.Router.Wide {
		if err := conn.Encode(msg); err != nil {
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
