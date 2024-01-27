package net

import (
	"context"
	"fmt"
	"time"

	"github.com/haiyanghan/tiangong/common/log"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
)

var (
	DefaultConnTimeout = 30 * time.Second
)

type TcpClient interface {
	Connect(handlerFunc ConnHandlerFunc) error
	Disconnect() error
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

	ctx := context.WithValue(t.ctx, common.TcpClientKey, t)
	if err = handlerFunc(ctx, t.conn); err != nil {
		log.Error("[TCP] connect closing....", err)
		_ = t.conn.Close()
		t.conn = nil
		return err
	}
	return err
}

func (t *tcpClientImpl) Write(buffer buf.Buffer) error {
	if t.conn != nil {
		if err := t.conn.ReadFrom(buffer); err != nil {
			return err
		}
		return nil
	}
	return errors.NewError("connect closed", nil)
}

func (t *tcpClientImpl) Disconnect() error {
	if t.conn != nil {
		_ = t.conn.Close()
		t.conn = nil
	}
	return nil
}

func NewTcpClient(host string, port int, ctx context.Context) TcpClient {
	return &tcpClientImpl{
		Host:    host,
		Port:    Port(port),
		Timeout: DefaultConnTimeout,
		ctx:     ctx,
	}
}
