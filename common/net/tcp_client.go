package net

import (
	"context"
	"fmt"
	"tiangong/common/buf"
	"tiangong/common/errors"
	"time"
)

var (
	DefaultConnTimeout = 30 * time.Second
)

type TcpClient interface {
	Connect(handlerFunc ConnHandlerFunc) error
	Write(buffer buf.Buffer) error
}

type tcpClientImpl struct {
	Host    string
	Port    Port
	Timeout time.Duration

	ctx  context.Context
	conn Conn
}

func (t *tcpClientImpl) Connect(handlerFunc ConnHandlerFunc) (err error) {
	if handlerFunc == nil {
		return errors.NewError("params handlerFunc Not be nil", nil)
	}
	t.conn, err = Dial("tcp", fmt.Sprintf("%s:%s", t.Host, t.Port.String()))
	if err != nil {
		return err
	}
	if err = handlerFunc(t.ctx, t.conn); err != nil {
		t.conn = nil
		_ = t.conn.Close()
		return err
	}
	return err
}

func (t *tcpClientImpl) Write(buffer buf.Buffer) error {
	if t.conn != nil {
		if err := t.conn.ReadFrom(buffer); err != nil {
			return err
		}
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
