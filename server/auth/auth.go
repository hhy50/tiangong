package auth

import (
	"google.golang.org/protobuf/proto"
	"net"
	"tiangong/common/errors"
	"tiangong/common/log"
	"tiangong/kernel/transport/protocol"
	"time"
)

var (
	Key     string
	TimeOut = 15 * time.Second
)

func Authentication(conn net.Conn) (proto.Message, error) {
	var err error
	if err = conn.SetDeadline(time.Now().Add(TimeOut)); err != nil {
		return nil, errors.NewError("Auth fail, SetDeadline error", err)
	}
	complete := func(status protocol.AuthStatus) {
		response := protocol.NewAuthResponse(status)
		var res []byte
		if res, err = response.Marshal(); err != nil {
			_ = conn.Close()
		}
		if _, err = conn.Write(res); err != nil {
			_ = conn.Close()
		}
	}

	var header *protocol.AuthHeader
	if header, err = protocol.DecodeAuthHeader(conn); err != nil {
		return nil, err
	}

	var body proto.Message = nil
	switch header.Type {
	case protocol.AuthClient:
		body = &protocol.ClientAuth{}
		break
	case protocol.AuthSession:
		body = &protocol.SessionAuth{}
		break
	default:
		return nil, errors.NewError("Unsupport AuthType: ["+string(header.Type)+"]", nil)
	}
	if err = protocol.DecodeProtoMessage(conn, int(header.Len), body); err != nil {
		return nil, errors.NewError("Auth fail, DecodeAuthBody error", err)
	}

	// 验证服务端key的有效性
	switch body.(type) {
	case *protocol.ClientAuth:
		clientAuth := body.(*protocol.ClientAuth)
		if clientAuth.Key != Key {
			complete(protocol.AuthFail)
			return nil, errors.NewError("Auth fail, client key not match", nil)
		}
		log.Info("New client join.")
		break
	case *protocol.SessionAuth:
		sessionAuth := body.(*protocol.SessionAuth)
		if err = Verification(sessionAuth.Token); err != nil {
			complete(protocol.AuthFail)
			return nil, errors.NewError("Auth fail", err)
		}
		log.Info("New session connected. token=%s, subHost=%s", sessionAuth.Token, sessionAuth.SubHost)
		break
	}

	complete(protocol.AuthSuccess)
	return body, nil
}
