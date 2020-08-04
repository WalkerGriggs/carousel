package identity

// Identity represnts the necessary information to authenticate with a Network.
// See RFC 2812 § 3.1
type Identity struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Realname string `json:"realname"`
	Password string `json:"password"`
}

func (i *Identity) CanAuthenticate() bool {
	return i.Username != "" && i.Password != ""
}

func (i *Identity) Wait() {
	for {
		if i.CanAuthenticate() {
			return
		}
	}
}
