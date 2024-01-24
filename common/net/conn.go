package net

import (
	"net"
	"tiangong/common/buf"
	"tiangong/common/log"
)

type Conn interface {
	net.Conn
	ReadFrom(buffer buf.Buffer) error
	Name() string
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
	log.Debug("write bytes [%x] --> [%s]", bytes, c.Name())
	if _, err = c.Write(bytes); err != nil {
		return err
	}
	return nil
}
