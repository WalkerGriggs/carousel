package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

type User struct {
	Username string
	Network  Network
	Router   Router
}

func (u User) ReadWrite() {
	conn, err := irc.Dial(u.Network.URI.Format())
	if err != nil {
		log.Fatal(err)
	}

	go u.Network.Identify(conn)

	go u.Read(conn)
	go u.Write(conn)
}

// (User -> Bouncer) <-> Client
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

// (User <- Bouncer) <-> Client
func (u User) Write(conn *irc.Conn) {
	for msg := range u.Router.Wide {
		if err := conn.Encode(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func Pong(conn *irc.Conn, msg *irc.Message) {
	conn.Encode(&irc.Message{
		Command: "PONG",
		Params:  msg.Params,
	})
}
