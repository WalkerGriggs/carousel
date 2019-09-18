package main

import (
	"log"

	"gopkg.in/sorcix/irc.v2"
)

// Network represents an IRC network. Each network has a URI, and, because Users
// own the Network object, each Network stores the User's Identity as well.
type Network struct {
	Name  string   `json:"name"`
	URI   URI      `json:"uri"`
	Ident Identity `json:"ident"`
}

// Identity represnts the necessary information to authenticate with a Network.
// See RFC 2812 § 3.1
type Identity struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
}

// Identify handles connection registration for each user.
// Again, see RFC 2812 § 3.1
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
