package server

import (
	"github.com/walkergriggs/carousel/pkg/user"
)

type Config struct {
	URI             string
	Users           []*user.User
	SSLEnabled      bool
	CertificatePath string
}
