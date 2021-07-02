package server

import (
	"crypto/tls"
	"net"

	"github.com/walkergriggs/carousel/pkg/crypto/ssl"
)

// listener returns a new network listener at the given address. This listener
// is not SSL encrypted.
func (s Server) listener() (net.Listener, error) {
	if s.config.SSLEnabled {
		return s.tlsListener()
	}

	return net.Listen("tcp", s.config.URI)
}

// tlsListener returns a network listener at the given address. This listener
// is ssl encryped using the configued certificate.
func (s Server) tlsListener() (net.Listener, error) {
	cert, err := ssl.LoadPem(s.config.CertificatePath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{*cert}}
	return tls.Listen("tcp", s.config.URI, config)
}
