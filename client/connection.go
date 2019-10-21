package client

import (
	"strings"

	"gopkg.in/sorcix/irc.v2"
)

func (c Client) Send(msg *irc.Message) error {
	_, err := c.Connection.Write([]byte(msg.String() + "\n"))
	return err
}

func (c Client) Receive() (*irc.Message, error) {
	msg, err := c.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	msg = strings.TrimSpace(string(msg))
	return irc.ParseMessage(msg), nil
}

func (c Client) BatchSend(messages []*irc.Message) error {
	for _, msg := range messages {
		if err := c.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
