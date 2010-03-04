package twitterstream

import "net"

type Client struct {
	Username string
	Password string
	conn *net.Conn
}

// Follow a list of user ids
func (c Client) Follow ([]string) {
}
