package carousel

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/walkergriggs/carousel/pkg/crypto/phash"
)

func (s *HTTPServer) usersRequest(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	switch req.Method {
	case "POST", "PUT":
		return s.userUpdate(resp, req)
	default:
		return nil, fmt.Errorf("ErrInvalidMethod")
	}
}

type (
	UserUpdateRequest struct {
		Username string
		Password string
	}

	UserUpdateResponse struct {
		Username string
	}
)

func (s *HTTPServer) userUpdate(resp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var out UserUpdateRequest
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&out); err != nil {
		return nil, err
	}

	pass, err := phash.Hash(out.Password)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(&UserConfig{
		Username: out.Username,
		Password: pass,
		Networks: make([]*Network, 0),
	})

	if err != nil {
		return nil, err
	}

	s.server.users = append(s.server.users, user)

	return &UserUpdateResponse{
		Username: user.Username,
	}, nil
}
