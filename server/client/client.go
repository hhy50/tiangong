package client

import (
	"net"
	"tiangong/client"
	"tiangong/server/auth"
	"tiangong/server/session"
)

type Client struct {
	Name string
}

type Auth struct {
}

func NewClient(name string, conn net.Conn) Client {

	return Client{}
}

func ConnHandler(conn net.Conn) {
	close := func(conn net.Conn) {
		_ = conn.Close()
	}

	user, err := auth.Authentication(conn)
	if err != nil {
		close(conn)
	}

	switch user.(type) {
	case client.Client:
		break
	case session.Session:
		break
	}
}
