package auth

import (
	"google.golang.org/protobuf/proto"
	"net"
	"tiangong/common/errors"
	"tiangong/kernel/transport/protocol"
	"tiangong/server"
	"time"
)

const (
	Client_Auth  protocol.AuthType = 1
	Session_Auth protocol.AuthType = 2
)

var (
	TimeOut = 15 * time.Second
)

func Authentication(conn net.Conn) (proto.Message, error) {
	if err := conn.SetDeadline(time.Now().Add(TimeOut)); err != nil {
		return nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	header, err := protocol.DecodeAuthHeader(conn)
	if err != nil {
		return nil, err
	}

	var body proto.Message = nil
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
		return nil, errors.NewError("Auth fail, DecodeAuthBody error", err)
	}

	// 验证服务端key的有效性
	switch body.(type) {
	case *protocol.ClientAuth:
		clientAuth := body.(*protocol.ClientAuth)
		if clientAuth.Key != server.Key {
			return nil, errors.NewError("Auth fail, client key not match", nil)
		}
		return body, nil
	case *protocol.SessionAuth:
		sessionAuth := body.(*protocol.SessionAuth)
		if err = Verification(sessionAuth.Token); err != nil {
			return nil, errors.NewError("Auth fail", err)
		}
	}
	return nil, nil
}
