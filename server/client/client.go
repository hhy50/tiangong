package client

import (
	"net"
	tgNet "tiangong/common/net"
	"tiangong/server"
)

type Client struct {
	Name string
	Host tgNet.IpAddress

	cli  server.Cli
	conn net.Conn
}
