package net

import (
	"net"
	"tiangong/common/buf"
)

type Conn interface {
	net.Conn
	ReadFrom(buffer buf.Buffer) error
}

type ConnWrap struct {
	net.Conn
}

func (c ConnWrap) ReadFrom(buffer buf.Buffer) error {
	bytes, err := buf.ReadAll(buffer)
	if err != nil {
		return err
	}
	if _, err = c.Write(bytes); err != nil {
		return err
	}
	return nil
}
