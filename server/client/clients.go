package client

import (
	"tiangong/common/errors"
	"tiangong/common/lock"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
)

var (
	// Clients with Router feature
	Clients     = make(map[net.IpAddress]*Client, 128)
	ClientNames = make(map[string]*Client, 128)
	Lock        = lock.NewLock()
)

func RegistClient(c *Client) error {
	Lock.Lock()
	defer Lock.Unlock()

	if _, f := Clients[c.Internal]; f {
		return errors.NewError("Unable to add existing client, name: "+c.Internal.String(), nil)
	}
	if _, f := ClientNames[c.Name]; f {
		return errors.NewError("Unable to add existing client, name: "+c.Name, nil)
	}
	Clients[c.Internal] = c
	ClientNames[c.Name] = c
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
