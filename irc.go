package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

func Connect(config Config, stob chan string, btoc chan *irc.Message) {
	conn, err := irc.Dial(config.URI.FormatURI())
	if err != nil {
		log.Fatal(err)
	}

	go Read(conn, stob)
	go Write(conn, btoc)

	Identify(conn, config)
}

func Write(conn *irc.Conn, btoc <-chan *irc.Message) {
	for msg := range btoc {
		if err := conn.Encode(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func Read(conn *irc.Conn, stob chan<- string) {
	for {
		msg, err := conn.Decode()
		if err != nil {
			log.Fatal(err)
		}

		if msg.Command == "PING" {
			Pong(conn, msg)
		}

		stob <- msg.String()
	}
}

func Identify(conn *irc.Conn, config Config) {
	var messages []*irc.Message
	messages = append(messages, &irc.Message{
		Command: irc.USER,
		Params:  []string{config.Nick, "0", "*", config.Nick},
	})

	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params:  []string{config.Nick},
	})

	BatchSend(messages, conn)
}

func BatchSend(messages []*irc.Message, conn *irc.Conn) {
	for _, msg := range messages {
		if err := conn.Encode(msg); err != nil {
			log.Fatal("Err: %s \n%s\n", err, msg)
		}
	}
}

func Pong(conn *irc.Conn, msg *irc.Message) {
	conn.Encode(&irc.Message{
		Command: "PONG",
		Params:  msg.Params,
	})
}
