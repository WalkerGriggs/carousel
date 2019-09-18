package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	URI      URI
	Listener net.Listener
}

type URI struct {
	Address string
	Port    int
	Type    string
}

func (uri URI) Format() string {
	return fmt.Sprintf("%s:%d", uri.Address, uri.Port)
}

func (s Server) Serve() {
	l, err := net.Listen("tcp", s.URI.Format())
	if err != nil {
		log.Fatal(err)
	}

	s.Listener = l
	defer l.Close()

	s.acceptLoop()
}

func (s Server) acceptLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go accept(conn)
	}
}

func accept(conn net.Conn) {
	connection := NewConnection(conn)
	router := NewRouter()

	ident := connection.decodeIdent(router)

	network := Network{
		Name: "Freenode",
		URI: URI{
			Address: "chat.freenode.net",
			Port:    6667,
			Type:    "tcp",
		},
		Ident: ident,
	}

	user := User{
		Username: ident.Username,
		Network:  network,
		Router:   router,
	}

	go connection.ReadWrite(user.Router)
	go user.ReadWrite()
}
