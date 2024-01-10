package client

import (
	"net"
	"tiangong/common/errors"
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

func NewClient(cli server.Cli, conn net.Conn) Client {
	return Client{
		Name: cli.Name,
		Host: cli.Internal,
		cli:  cli,
		conn: conn,
	}
}
