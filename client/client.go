package client

import (
	"tiangong/common/conf"
	"tiangong/common/net"
	"time"
)

type Client struct {
	cnf       Config
	tcpClient net.TcpClient
}

func (s *Client) Start() {

}

func handshake(conn net.Conn) {

}

// NewClient by specify a config file
func NewClient(cp string) (*Client, error) {
	c := Config{}
	if err := conf.LoadConfig(cp, c, defaultValue); err != nil {
		return nil, err
	}

	tc := net.TcpClient{
		Host:    c.ServerHost,
		Port:    c.ServerPort,
		Timeout: 30 * time.Second,
	}
	return &Client{c, tc}, nil
}
