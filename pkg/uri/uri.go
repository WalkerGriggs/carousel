package uri

import (
	"fmt"
)

type Options struct {
	Host string
	Port int
}

// URI is the basic information needed to address Networks and Servers. URI is
// not an exhaustive liste of all URI components, and will be extended in future
// implementations.
type URI struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func New(ops Options) (*URI, error) {
	return &URI{
		Host: ops.Host,
		Port: ops.Port,
	}, nil
}

func (uri URI) String() string {
	return fmt.Sprintf("%s:%d", uri.Host, uri.Port)
}
