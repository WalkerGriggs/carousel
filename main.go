package main

import (
	"fmt"

	"gopkg.in/sorcix/irc.v2"
)

type URI struct {
	Address string
	Port    int
	Type    string
}

type Config struct {
	URI  URI
	Nick string
}

func (uri URI) FormatURI() string {
	return fmt.Sprintf("%s:%d", uri.Address, uri.Port)
}

func PrintAll(messages <-chan string) {
	for msg := range messages {
		fmt.Println(msg)
	}
}

func main() {
	client := Config{
		URI: URI{
			Address: "irc.freenode.net",
			Port:    6667,
		},
		Nick: "rub2k",
	}

	stob := make(chan string)
	btoc := make(chan *irc.Message)

	go Connect(client, stob, btoc)

	bouncer := Bouncer{
		URI: URI{
			Address: "0.0.0.0",
			Port:    6667,
		},

		Stob: stob,
		Btoc: btoc,
	}

	bouncer.Serve()
}
