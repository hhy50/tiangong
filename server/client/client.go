package client

import "net"

type Client struct {
	Name string

}

func NewClient(name string, conn net.Conn) Client {

	return Client{

	}
}


func ConnHanlder(conn net.Conn) {

}