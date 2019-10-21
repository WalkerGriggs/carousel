package network

import (
	log "github.com/sirupsen/logrus"
)

func (n Network) LogEntry() *log.Entry {
	return log.WithFields(n.LogFields())
}

func (n Network) LogFields() log.Fields {
	return log.Fields{
		"Network": n.Name,
		"Host":    n.URI.String(),
		"User":    n.Ident.Username,
	}
}
