package net

import (
	"fmt"
	"tiangong/common/errors"
	"time"
)

type TcpClient struct {
	Host    string
	Port    Port
	Timeout time.Duration

	conn Conn
}

func (t *TcpClient) Connect(handlerFunc ConnHandlerFunc) error {
	if handlerFunc == nil {
		return errors.NewError("params handlerFunc Not be nil", nil)
	}
	conn, err := Dial("tcp", fmt.Sprintf("%s:%s", t.Host, t.Port.String()))
	if err != nil {
		return err
	}
	if err := handlerFunc(conn); err != nil {
		_ = conn.Close()
		return err
	}
	t.conn = conn
	return nil
}

func (t *TcpClient) Disconnect() {

}
