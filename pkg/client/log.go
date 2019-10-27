package client

import (
	log "github.com/sirupsen/logrus"
)

func (c *Client) LogEntry() *log.Entry {
	return log.WithFields(c.LogFields())
}

func (c *Client) LogFields() log.Fields {
	return log.Fields{
		"RemoteAddress": c.Connection.RemoteAddr().String,
	}
}
