package network

import (
	log "github.com/sirupsen/logrus"
)

func (n *Network) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"Network": n.Name,
		"Host":    n.URI.String(),
		"User":    n.Ident.Username,
	})
}
