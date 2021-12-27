package carousel

import (
	"gopkg.in/sorcix/irc.v2"
)

type Connection interface {
	Send(msg *irc.Message) error
	Receive() (*irc.Message, error)
	BatchSend(message []*irc.Message) error
}
