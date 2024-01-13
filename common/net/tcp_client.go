package net

import (
	"fmt"
	"net"
	"tiangong/common/errors"
	"time"
)

type TcpClient struct {
	Host    string
	Port    Port
	Timeout time.Duration
}

func (t *TcpClient) Conn(handlerFunc ConnHandlerFunc) error {
	if handlerFunc == nil {
		return errors.NewError("params handlerFunc Not be nil", nil)
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", t.Host, t.Port.String()))
	if err != nil {
		return err
	}
	return handlerFunc(conn)
}
