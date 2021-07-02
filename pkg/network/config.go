package network

import (
	"github.com/walkergriggs/carousel/pkg/identity"
)

type Config struct {
	Name  string
	URI   string
	Ident *identity.Identity
}
