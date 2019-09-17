package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"gopkg.in/sorcix/irc.v2"
)

type Bouncer struct {
	URI  URI
	Stob chan string
	Btoc chan *irc.Message
}

func (b Bouncer) Serve() {
	listener, err := net.Listen("tcp", b.URI.FormatURI())
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Listening on " + b.URI.FormatURI())
	acceptLoop(b, listener)
}

func acceptLoop(config Bouncer, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go readLoop(config, conn)
		go writeLoop(config, conn)
	}
}

func readLoop(config Bouncer, conn net.Conn) {
	fmt.Printf("Serving %s\n", conn.RemoteAddr().String())

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		temp := strings.TrimSpace(string(netData))

		if message := irc.ParseMessage(temp); message != nil {
			config.Btoc <- message
		}
	}

	conn.Close()
}

func writeLoop(config Bouncer, conn net.Conn) {
	for msg := range config.Stob {
		if _, err := conn.Write([]byte(string(msg) + "\n")); err != nil {
			log.Fatal(err)
		}
	}
}
