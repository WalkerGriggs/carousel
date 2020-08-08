package channel

type Channel struct {
	Name string `json:",omitempty"`
	Nicks []string `json:",omitempty"`
}

func New(name string) (*Channel, error) {
	return &Channel{
		Name: name,
		Nicks: make([]string, 1),
	}, nil
}

func (c *Channel) AddNicks(nicks []string) {
	for _, nick := range nicks {
		c.AddNick(nick)
	}
}

func (c *Channel) AddNick(nick string) {
	if !c.hasNick(nick) {
		c.Nicks = append(c.Nicks, nick)
	}
}

func (c *Channel) RemoveNick(nick string) {
	if i := c.nickIndex(nick); i != -1 {
		c.Nicks = append(c.Nicks[:i], c.Nicks[i+1:]...)
	}
}

func (c *Channel) hasNick(nick string) bool {
	return c.nickIndex(nick) != -1
}

func (c *Channel) nickIndex(nick string) int {
	for i, n := range c.Nicks {
		if n == nick {
			return i
		}
	}
	return -1
}
