package carousel

import (
	"fmt"
	"net/http"
)

func (s *HTTPServer) pongRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "GET":
		return s.pong(resp, req)
	default:
		return nil, fmt.Errorf("ErrInvalidMethod")
	}
}

func (s *HTTPServer) pong(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	return "pong", nil
}
