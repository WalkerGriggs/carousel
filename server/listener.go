package server

import (
	"crypto/tls"
	"net"

	"github.com/walkergriggs/carousel/crypto/ssl"
)

func (s Server) listener() (net.Listener, error) {
	if s.SSLEnabled {
		return s.tlsListener()
	}

	return net.Listen("tcp", s.URI.String())
}

func (s Server) tlsListener() (net.Listener, error) {
	cert, err := ssl.LoadPem(s.CertificatePath)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{Certificates: []tls.Certificate{*cert}}
	return tls.Listen("tcp", s.URI.String(), config)
}
