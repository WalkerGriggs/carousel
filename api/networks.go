package api

type Networks struct {
	client *Client
}

func (c *Client) Networks() *Networks {
	return &Networks{client: c}
}

type (
	Network struct {
		Name     string
		URI      string
		Ident    *Identity
		Channels []string
	}

	Identity struct {
		Username string
		Password string
		Realname string
		Nickname string
	}
)

type (
	NetworkCreateRequest struct {
		User    string
		Network *Network
	}

	NetworkCreateResponse struct {
		Name string
	}
)

func (n *Networks) Create(user string, network *Network) (*NetworkCreateResponse, error) {
	req := &NetworkCreateRequest{
		User:    user,
		Network: network,
	}

	var res NetworkCreateResponse
	if err := n.client.write("/v1/networks", req, &res, nil); err != nil {
		return nil, err
	}

	return &res, nil
}
