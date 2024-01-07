package client

import (
	"net"
	tgNet "tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"tiangong/server/auth"
)

type Client struct {
	Name string
	Host tgNet.IpAddress

	header *protocol.AuthHeader
	conn   net.Conn
}

func NewClient(name string, host tgNet.IpAddress, header *protocol.AuthHeader, conn net.Conn) Client {
	return Client{
		Name:   name,
		Host:   host,
		header: header,
		conn:   conn,
	}
}

func ConnHandler(conn net.Conn) {
	close := func() {
		_ = conn.Close()
	}

	user, err := auth.Authentication(conn)
	if err != nil {
		close()
	}

	switch user.(type) {
	//case client.Client:
	//	break
	//case session.Session:
	//	break
	}
}
