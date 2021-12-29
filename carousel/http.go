package carousel

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

type (
	HTTPWrapper func(resp http.ResponseWriter, req *http.Request) (interface{}, error)
	HTTPHandler func(resp http.ResponseWriter, req *http.Request)
)

type HTTPConfig struct {
	Advertise string
}

type HTTPServer struct {
	server     *Server
	addr       string
	mux        *http.ServeMux
	listener   net.Listener
	listenerCh chan struct{}
}

// NewHTTPServer configures a new multiplexer and listener and returns the
// running HTTPServer.
func NewHTTPServer(s *Server, c *HTTPConfig) (*HTTPServer, error) {
	mux := http.NewServeMux()

	addr, err := net.ResolveTCPAddr("tcp", c.Advertise)
	if err != nil {
		return nil, err
	}

	hostport := net.JoinHostPort(addr.IP.String(), strconv.Itoa(addr.Port))

	ln, err := net.Listen("tcp", hostport)
	if err != nil {
		return nil, fmt.Errorf("Failed to start HTTP listener: %v", err)
	}

	srv := &HTTPServer{
		server:     s,
		mux:        mux,
		listener:   ln,
		listenerCh: make(chan struct{}),
	}

	srv.registerHandlers()

	httpServer := http.Server{
		Addr:    srv.addr,
		Handler: srv.mux,
	}

	go func() {
		defer close(srv.listenerCh)
		httpServer.Serve(ln)
	}()

	return srv, nil
}

// registerHandlers maps each handler to an endpoint on the HTTPServer's
// multiplexer.
func (s *HTTPServer) registerHandlers() {
	s.mux.HandleFunc("/v1/ping", s.wrap(s.pongRequest))
	s.mux.HandleFunc("/v1/users", s.wrap(s.usersRequest))
	s.mux.HandleFunc("/v1/networks", s.wrap(s.networksRequest))
}

// wrap wraps the handler function with some quality-of-life improvements. It
// returns a net/http ServeMux compliant handler function.
func (s *HTTPServer) wrap(wrapper HTTPWrapper) HTTPHandler {
	return func(resp http.ResponseWriter, req *http.Request) {
		obj, err := wrapper(resp, req)

	HAS_ERR:
		if err != nil {
			resp.WriteHeader(500)
			resp.Write([]byte(err.Error()))
			return
		}

		if obj != nil {
			bytes, err := json.Marshal(obj)
			if err != nil {
				goto HAS_ERR
			}

			resp.Header().Set("Content-Type", "application/json")
			resp.Write(bytes)
		}
	}
}
