package client

import (
	"fmt"
	"net"
	"tiangong/common/buf"
)

type Client struct {
	host   string
	port   int
	conn   net.Conn
	buffer buf.Buffer
}

func (s *Client) connect() {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		panic(err)
	}
	s.conn = conn

	go func(c net.Conn) {
		buf := make([]byte, 1024)
		for {
			len, err := c.Read(buf)
			if len == 0 {
				return
			}
			if err != nil {
				return
			}
			//msg := string(buf[:len-1])
		}
	}(s.conn)
}

func (s *Client) Write(msg []byte) {
	s.conn.Write(msg)
}

func NewClient(host string, port int) *Client {
	client := &Client{
		host: host,
		port: port,
	}
	client.connect()
	return client
}
