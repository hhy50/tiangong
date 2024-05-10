package net

import (
	"net"

	"github.com/haiyanghan/tiangong/common/buf"
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
