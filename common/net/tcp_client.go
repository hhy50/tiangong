package net

import (
	"net"
)

type ConnHandler func(net.Conn)

type TcpClient struct {
	Host string
	Port string
}

func (s *TcpClient) Conn() {

}


