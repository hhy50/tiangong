package client

import (
	"tiangong/common/errors"
	"tiangong/common/net"
	"tiangong/server"
)

var (
	Clients = make(map[string]*Client, 128)
)

func AddClient(c *Client) error {
	name := c.Name
	if _, f := Clients[name]; f {
		return errors.NewError("Unable to add existing client, name: "+name, nil)
	}
	Clients[name] = c
	return nil
}

func NewClient(internalIP net.IpAddress, cli server.Cli, conn net.Conn) Client {
	return Client{
		Name: cli.Name,
		Host: internalIP,
		cli:  cli,
		conn: conn,
	}
}
