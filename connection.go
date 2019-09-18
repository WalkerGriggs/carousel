package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"gopkg.in/sorcix/irc.v2"
)

type Connection struct {
	Conn    net.Conn
	Encoder *irc.Encoder
	Decoder *irc.Decoder
}

func NewConnection(conn net.Conn) Connection {
	return Connection{
		Conn:    conn,
		Encoder: irc.NewEncoder(conn),
		Decoder: irc.NewDecoder(conn),
	}
}

// TODO: Fix assumption that PASS, NICK, USER arrive in order
func (c Connection) decode(command string, router Router) *irc.Message {
	for {
		message, err := c.Decoder.Decode()
		if err != nil {
			log.Fatal(err)
		}

		router.Wide <- message

		if message.Command == command {
			return message
		}
	}
}

func (c Connection) decodeIdent(router Router) Identity {
	messages := make(map[string]*irc.Message)
	required_commands := []string{"USER", "NICK", "PASS"}

	for {
		message, err := c.Decoder.Decode()
		if err != nil {
			log.Fatal(err)
		}

		messages[strings.ToUpper(message.Command)] = message

		if containsAll(messages, required_commands) {
			break
		}
	}

	return Identity{
		Nickname: messages["NICK"].Params[0],
		Username: messages["USER"].Params[0],
		Realname: messages["USER"].Params[3],
		Password: messages["PASS"].Params[0],
	}
}

func (c Connection) ReadWrite(router Router) {
	go c.Read(router)
	go c.Write(router)
}

// User <-> (Bouncer <- Client)
func (c Connection) Read(router Router) {
	reader := bufio.NewReader(c.Conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		msg = strings.TrimSpace(string(msg))

		if irc_message := irc.ParseMessage(msg); irc_message != nil {
			router.Wide <- irc_message
		}
	}
}

// User <-> (Bouncer -> Client)
func (c Connection) Write(router Router) {
	for {
		msg := <-router.Local
		if _, err := c.Conn.Write([]byte(msg.String() + "\n")); err != nil {
			log.Fatal(err)
		}
	}
}

func containsAll(messages map[string]*irc.Message, required_commands []string) bool {
	for _, command := range required_commands {
		if _, ok := messages[command]; !ok {
			return false
		}
	}

	return true
}
