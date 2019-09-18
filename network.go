package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

type Network struct {
	Name  string   `json:"name"`
	URI   URI      `json:"uri"`
	Ident Identity `json:"ident"`
}

type Identity struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
}

func (net Network) Identify(conn *irc.Conn) {
	var messages []*irc.Message

	messages = append(messages, &irc.Message{
		Command: irc.PASS,
		Params:  []string{net.Ident.Password},
	})

	messages = append(messages, &irc.Message{
		Command: irc.NICK,
		Params:  []string{net.Ident.Nickname},
	})

	messages = append(messages, &irc.Message{
		Command: irc.USER,
		Params:  []string{net.Ident.Username, "0", "*", net.Ident.Realname},
	})

	batchSend(messages, conn)
}

func batchSend(messages []*irc.Message, conn *irc.Conn) {
	for _, msg := range messages {
		if err := conn.Encode(msg); err != nil {
			log.Fatal("Err: %s \n%s\n", err, msg)
		}
	}
}
