package server

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/walkergriggs/carousel/pkg/client"
	"github.com/walkergriggs/carousel/pkg/router"
	"github.com/walkergriggs/carousel/pkg/uri"
	"github.com/walkergriggs/carousel/pkg/user"
)

type Options struct {
	URI             uri.URI
	Users           []*user.User
	SSLEnabled      bool
	CertificatePath string
}

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	URI             uri.URI      `json:"uri"`
	Users           []*user.User `json:"users"`
	SSLEnabled      bool         `json:"sslEnabled"`
	CertificatePath string       `json:"certificatePath"`
	Listener        net.Listener `json:",omitempty"`
}

func New(opts Options) (*Server, error) {
	return &Server{
		URI:             opts.URI,
		Users:           opts.Users,
		SSLEnabled:      opts.SSLEnabled,
		CertificatePath: opts.CertificatePath,
	}, nil
}

// Serve attaches a tcp listener to the specificed URI, and starts the main
// event loop. Serve blocks for the lifetime of the parent process and should
// only return if the TCP listener closes or errors (even if there are no active
// connections).
func (s Server) Serve() {
	l, err := s.listener()
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	ctx := context.Background()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
		}

		go s.acceptConnection(ctx, conn)
	}
}

// acceptConnection establishes a new connection with the accepted TCP client, and spawns
// the concurrent processess responsible to message passing between the IRC
// network and user. Each accepted connection has it's own router and associated
// user, so accept should only return when the user disconnects, or does not
// authenticate.
func (s Server) acceptConnection(ctx context.Context, conn net.Conn) {
	ctx, cancel := context.WithCancel(ctx)
	c, _ := client.New(client.Options{
		Connection: conn,
		Ctx:        ctx,
		Cancel:     cancel,
	})
	c.Listen()

	// authorize is a blocking function, and will not return until the user has
	// been authroized or (todo) timeout reached
	u, err := s.registerClient(c)
	if err != nil {
		c.LogEntry().WithError(err).Error("Failed to authorize client.")
		c.Disconnect()
		return
	}

	router := router.Router{
		Client:  c,
		Network: u.Network,
		Ctx:     ctx,
	}

	go u.Network.Listen()
	go router.AttachClient()
	go router.Route()
}
