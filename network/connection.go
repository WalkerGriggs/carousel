package network

import (
	"gopkg.in/sorcix/irc.v2"
)

func (n Network) Send(msg *irc.Message) error {
	return n.Connection.Encode(msg)
}

func (n Network) Receive() (*irc.Message, error) {
	return n.Connection.Decode()
}

func (n Network) BatchSend(messages []*irc.Message) error {
	for _, msg := range messages {
		if err := n.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
