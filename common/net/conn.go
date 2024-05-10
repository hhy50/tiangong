package net

import (
	"net"
	"time"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
)

const (
	ConnValName = "net_conn"
)

type Conn interface {
	net.Conn
	ReadFrom(buffer buf.Buffer) error
}

type ConnWrap struct {
	net.Conn
}

func (c ConnWrap) Name() string {
	return c.RemoteAddr().String()
}

func (c ConnWrap) ReadFrom(buffer buf.Buffer) error {
	if err := c.Conn.SetWriteDeadline(time.Now().Add(15*time.Second)); err != nil {
		return errors.NewError("SetWriteDeadline err", err)
	}

	bytes, err := buf.ReadAll(buffer)
	if err != nil {
		return err
	}
	log.Debug("Write %d size bytes [%x] --> [%s]", len(bytes), bytes, c.Name())
	if _, err = c.Write(bytes); err != nil {
		return err
	}
	return nil
}
