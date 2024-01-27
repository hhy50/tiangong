package net

import (
	"net"
)

type OpError = net.OpError

func Dial(network, address string) (Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return ConnWrap{conn}, nil
}
