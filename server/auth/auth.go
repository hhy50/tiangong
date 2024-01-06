package auth

import (
	"net"
	"tiangong/common/errors"
	"tiangong/kernel/transport/protocol"
	"time"
)

const (
	Client_Auth  protocol.AuthType = 1
	Session_Auth protocol.AuthType = 2
)

var (
	TimeOut = 15 * time.Second
)

func Authentication(conn net.Conn) (interface{}, error) {
	if err := conn.SetDeadline(time.Now().Add(TimeOut)); err != nil {
		return nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	header, err := protocol.DecodeAuthHeader(conn)
	if err != nil {
		return nil, err
	}

	var body interface{} = nil
	switch header.Type {
	case Client_Auth:
		body, err = protocol.DecodeClientAuthBody(conn, header.Len)
		break
	case Session_Auth:
		break
	default:
		return nil, errors.NewError("Unsupport AuthType: ["+string(header.Type)+"]", nil)
	}
	if err != nil || body == nil {
		return nil, err
	}

	// TODO 验证服务端key的有效性
	return nil, nil
}
