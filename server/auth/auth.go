package auth

import (
	"tiangong/common"
	"tiangong/common/buf"
	"tiangong/common/errors"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport"
	"tiangong/kernel/transport/protocol"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	TimeOut = 15 * time.Second
)

func Authentication(key string, conn net.Conn) (*protocol.AuthHeader, proto.Message, error) {
	buffer := buf.NewBuffer(256)
	defer buffer.Release()

	if err := conn.SetDeadline(time.Now().Add(TimeOut)); err != nil {
		return nil, nil, errors.NewError("Auth fail, SetDeadline error", err)
	}
	complete := func(status protocol.AuthStatus) {
		response := protocol.NewAuthResponse(status)
		if err := response.WriteTo(buffer); err != nil {
			_ = conn.Close()
		}

		if err := conn.ReadFrom(buffer); err != nil {
			_ = conn.Close()
		}
	}

	var header protocol.AuthHeader
	if err := protocol.DecodeAuthHeader(conn, &header); err != nil {
		return nil, nil, err
	}

	var body proto.Message = nil
	switch header.Type {
	case protocol.AuthClient:
		body = &protocol.ClientAuth{}
	case protocol.AuthSession:
		body = &protocol.SessionAuth{}
	default:
		return nil, nil, errors.NewError("Unsupport AuthType: ["+string(header.Type)+"]", nil)
	}
	if err := transport.DecodeProtoMessage(conn, int(header.Len), body); err != nil {
		return nil, nil, errors.NewError("Auth fail, DecodeAuthBody error", err)
	}

	// 验证服务端key的有效性
	switch body.(type) {
	case *protocol.ClientAuth:
		clientAuth := body.(*protocol.ClientAuth)
		if common.IsNotEmpty(key) && clientAuth.Key != key {
			complete(protocol.AuthFail)
			return nil, nil, errors.NewError("Auth fail, client key not match", nil)
		}
		log.Info("New client join. ")
	case *protocol.SessionAuth:
		sessionAuth := body.(*protocol.SessionAuth)
		if err := Verification(sessionAuth.Token); err != nil {
			complete(protocol.AuthFail)
			return nil, nil, errors.NewError("Auth fail", err)
		}
		log.Info("New session connected. token=%s, subHost=%s", sessionAuth.Token, sessionAuth.SubHost)
	default:
		return nil, nil, errors.NewError("Not support auth type", nil)
	}

	complete(protocol.AuthSuccess)
	return &header, body, nil
}
