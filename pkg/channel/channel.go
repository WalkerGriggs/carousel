package channel

type Channel struct {
	Name string `json:",omitempty"`
}

func New(name string) (*Channel, error) {
	return &Channel{
		Name: name,
	}, nil
}
