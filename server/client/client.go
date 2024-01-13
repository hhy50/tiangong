package client

import (
	"net"
	tgNet "tiangong/common/net"
	"tiangong/kernel/transport/protocol"
)

type Client struct {
	Name string
	Host tgNet.IpAddress

	cli  *protocol.ClientAuth
	conn net.Conn
}
