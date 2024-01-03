package auth

import (
	"fmt"
	"net"
	"tiangong/common/errors"
	"time"
)

var (
	TimeOut = 15 * time.Second
	AuthLen = 8
)

func Authentication(conn net.Conn) (interface{}, error) {
	if err := conn.SetDeadline(time.Now().Add(TimeOut)); err != nil {
		return nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	bytes := make([]byte, AuthLen)
	if n, err := conn.Read(bytes); err != nil || n != AuthLen {
		return 0, errors.NewError(fmt.Sprintf("Auth fial, expect read %d bytes, actuality read %d bytes", AuthLen, n), err)
	}
	return nil, nil
}
