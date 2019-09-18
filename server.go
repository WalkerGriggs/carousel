package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	URI   URI    `json:"uri"`
	Users []User `json:"users"`

	Listener net.Listener
}

type URI struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
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

		go s.accept(conn)
	}
}

func (s Server) accept(conn net.Conn) {
	connection := NewConnection(conn)
	router := NewRouter()

	username := connection.decodeIdent(router)
	fmt.Println(username)

	user := *getUser(username, s.Users)
	user.Router = router

	go connection.ReadWrite(user.Router)
	go user.ReadWrite()
}

func getUser(username string, users []User) *User {
	for _, user := range users {
		if username == user.Username {
			return &user
		}
	}

	return nil
}
