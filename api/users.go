package api

// Users wraps the client and is used for task-specific endpoints
type Users struct {
	client *Client
}

// FSMs wraps the client for task-specific endpoints
func (c *Client) Users() *Users {
	return &Users{client: c}
}

type User struct {
	Username string
	Password string
}

type (
	// UserCreateRequest is used to serialize a Define request
	UserCreateRequest struct {
		Username string
		Password string
	}

	// UserCreateResponse is used to serialize a Define response
	UserCreateResponse struct {
		Username string
	}
)

func (u *Users) Create(user *User) (*UserCreateResponse, error) {
	req := &UserCreateRequest{
		Username: user.Username,
		Password: user.Password,
	}

	var res UserCreateResponse
	if err := u.client.write("/v1/users", req, &res, nil); err != nil {
		return nil, err
	}

	return &res, nil
}
