package net

import (
	"context"
	"fmt"
	"tiangong/common/errors"
	"time"
)

var (
	DefaultConnTimeout = 30 * time.Second
)

type TcpClient interface {
	Connect(handlerFunc ConnHandlerFunc) error
}

type tcpClientImpl struct {
	Host    string
	Port    Port
	Timeout time.Duration

	ctx  context.Context
	conn Conn
}

func (t *tcpClientImpl) Connect(handlerFunc ConnHandlerFunc) error {
	if handlerFunc == nil {
		return errors.NewError("params handlerFunc Not be nil", nil)
	}
	conn, err := Dial("tcp", fmt.Sprintf("%s:%s", t.Host, t.Port.String()))
	if err != nil {
		return err
	}
	t.conn = conn
	if err := handlerFunc(t.ctx, conn); err != nil {
		_ = conn.Close()
		return err
	}
	return nil
}

func (t *tcpClientImpl) Disconnect() {

}

func NewTcpClient(host string, port int, ctx context.Context) TcpClient {
	return &tcpClientImpl{
		Host:    host,
		Port:    Port(port),
		Timeout: DefaultConnTimeout,
		ctx:     ctx,
	}
}
