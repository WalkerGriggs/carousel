package carousel

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/sorcix/irc.v2"
)

func (s *HTTPServer) networksRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "POST", "PUT":
		return s.networkUpdate(resp, req)
	default:
		return nil, fmt.Errorf("ErrInvalidMethod")
	}
}

type (
	NetworkUpdateRequest struct {
		User    string
		Network *Network
	}

	NetworkUpdateResponse struct {
		Name string
	}
)

func (s *HTTPServer) networkUpdate(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var out NetworkUpdateRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}

	// TODO use NewNetwork instead
	//      add Ident completion (?)
	out.Network.Buffer = make(chan *irc.Message)

	for _, user := range s.server.users {
		if user.Username == out.User {
			user.Networks = append(user.Networks, out.Network)

			return &NetworkUpdateResponse{
				Name: out.Network.Name,
			}, nil
		}
	}

	return nil, fmt.Errorf("Cannot find user %s", out.User)
}
