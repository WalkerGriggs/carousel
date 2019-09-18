package main

import (
	"gopkg.in/sorcix/irc.v2"
)

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
