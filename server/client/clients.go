package client

import (
	"tiangong/common/errors"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
)

var (
	// Clients with Router feature
	Clients = make(map[net.IpAddress]*Client, 128)
)

func RegistClient(c *Client) error {
	name := c.Internal
	if _, f := Clients[name]; f {
		return errors.NewError("Unable to add existing client, name: "+name.String(), nil)
	}
	Clients[name] = c
	return nil
}

func NewClient(internalIP net.IpAddress, cli *protocol.ClientAuth, conn net.Conn) Client {
	return Client{
		Name:     cli.Name,
		Internal: internalIP,
		auth:     cli,
		conn:     conn,
	}
}
