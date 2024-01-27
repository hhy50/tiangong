package auth

import (
	"fmt"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport"
	"github.com/haiyanghan/tiangong/transport/protocol"
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
		_ = buffer.Clear()
		log.Debug("write auth response body, status:[%d]", status)
		response := protocol.NewAuthResponse(status)
		if err := response.WriteTo(buffer); err != nil {
			log.Warn("write to auth response error", err)
			return
		}

		if err := conn.ReadFrom(buffer); err != nil {
			log.Warn("write to auth response error", err)
			return
		}
	}

	if n, err := buffer.Write(conn, protocol.AuthHeaderLen); err != nil || n != protocol.AuthHeaderLen {
		return nil, nil, errors.NewError(
			fmt.Sprintf("read bytes from connect too short, should minnum read %d bytes actual reading %d bytes",
				protocol.AuthHeaderLen, n), err)
	}

	var header protocol.AuthHeader
	if err := protocol.DecodeAuthHeader(buffer, &header); err != nil {
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
		complete(protocol.AuthSuccess)
		log.Info("New client join. name: [%s], internal:[%s]", clientAuth.Name, net.ValueOf(clientAuth.Internal).String())
	case *protocol.SessionAuth:
		sessionAuth := body.(*protocol.SessionAuth)
		if err := Verification(sessionAuth.Token); err != nil {
			complete(protocol.AuthFail)
			return nil, nil, errors.NewError("Auth fail", err)
		}
		complete(protocol.AuthSuccess)
		log.Info("New session connected. token=%s, subHost=%s", sessionAuth.Token, sessionAuth.SubHost)
	default:
		return nil, nil, errors.NewError("Not support auth type", nil)
	}
	return &header, body, nil
}
