package user

import (
	"github.com/walkergriggs/carousel/pkg/network"
)

type Config struct {
	Username string
	Password string
	Networks []*network.Network
}
