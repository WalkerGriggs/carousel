package carousel

type Channel struct {
	Name  string   `json:",omitempty"`
	Nicks []string `json:",omitempty"`
}

func NewChannel(name string) (*Channel, error) {
	return &Channel{
		Name:  name,
		Nicks: make([]string, 1),
	}, nil
}

// AddNicks adds a list of nicks to the channel.
func (c *Channel) AddNicks(nicks []string) {
	for _, nick := range nicks {
		c.AddNick(nick)
	}
}

// AddNick adds a given nick to the channel.
func (c *Channel) AddNick(nick string) {
	if !c.hasNick(nick) {
		c.Nicks = append(c.Nicks, nick)
	}
}

// RemoveNick removes a given nick form the channel.
func (c *Channel) RemoveNick(nick string) {
	if i := c.nickIndex(nick); i != -1 {
		c.Nicks = append(c.Nicks[:i], c.Nicks[i+1:]...)
	}
}

// hasNick is a predicate which checks if the channel already has a given nick
// recoreded. It returns a boolean.
func (c *Channel) hasNick(nick string) bool {
	return c.nickIndex(nick) != -1
}

// nickIndex returns an index for the given nick in the channels list of nicks.
// It returns -1 if the channel doesn't have the nick recorded.
func (c *Channel) nickIndex(nick string) int {
	for i, n := range c.Nicks {
		if n == nick {
			return i
		}
	}
	return -1
}
