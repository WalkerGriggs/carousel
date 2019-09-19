package main

// User represents the individual Users's account and network config. Currently,
// a user can only connect to a single network, and each user owns their own
// router to pass messages.
type User struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	Network  Network `json:"network"`

	Router Router
}
