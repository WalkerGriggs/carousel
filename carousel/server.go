package carousel

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/walkergriggs/carousel/pkg/rungroup"
)

type ServerConfig struct {
	URI             string
	Users           []*User
	SSLEnabled      bool
	CertificatePath string
}

// Server is the configuration for all of Carousel. It maintains a list of all
// Users, as well general server information (ie. URI).
type Server struct {
	config *ServerConfig
	users  []*User
}

func NewServer(config *ServerConfig) (*Server, error) {
	return &Server{
		config: config,
		users:  config.Users,
	}, nil
}

// Serve attaches a tcp listener to the specificed URI, and starts the main
// event loop. Serve blocks for the lifetime of the parent process and should
// only return if the TCP listener closes or errors (even if there are no active
// connections).
func (s Server) Serve() {
	log.Info("Listening at ", s.config.URI)

	l, err := s.listener()
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
		}

		go s.acceptConnection(conn)
	}
}

// acceptConnection establishes a new connection with the accepted TCP client,
// and spawns the concurrent processess responsible to message passing between
// the network and user.
func (s Server) acceptConnection(conn net.Conn) {
	clientGroup := rungroup.New(context.Background())

	c, _ := NewClient(ClientConfig{conn})
	clientGroup.Go(c.Listen)
	clientGroup.Go(c.Heartbeat)

	u, err := s.blockingAuthorizeClient(c)
	if err != nil {
		c.LogEntry().Error("Failed to authorize client.")
		c.Disconnect()
		return
	}

	network := u.NetworkOrDefault(c.Ident.ParsedNetwork())

	router := Router{
		Client:    c,
		Network:   network,
		ServerURI: s.config.URI,
	}

	go network.Listen()
	go router.attachClient()
	clientGroup.Go(router.Route)

	if err := clientGroup.Wait(); err != nil {
		log.Error(err)
	}
}
