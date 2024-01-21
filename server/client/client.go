package client

import (
	"tiangong/common/buf"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
)

type Client struct {
	Name     string
	Internal net.IpAddress

	auth *protocol.ClientAuth
	conn net.Conn
}

func (c *Client) Write(buffer buf.Buffer) error {
	return c.conn.ReadFrom(buffer)
}
