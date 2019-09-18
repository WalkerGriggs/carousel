package main

import (
	"gopkg.in/sorcix/irc.v2"
)

// Router maintains the two channels over which messages are passed between the
// User and Network. Named after WAN and LAN, the Local channel handles all
// traffic between the User and Router, and the Wide channel handles all traffic
// between the Router and Network
type Router struct {
	Local chan *irc.Message
	Wide  chan *irc.Message
}

func NewRouter() Router {
	return Router{
		Local: make(chan *irc.Message),
		Wide:  make(chan *irc.Message),
	}
}
