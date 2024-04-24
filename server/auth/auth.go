package auth

import (
	"fmt"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport"
	"github.com/haiyanghan/tiangong/transport/protocol"

	"google.golang.org/protobuf/proto"
)

var (
	Timout   = 15 * time.Second
	complete = func(conn net.Conn, status protocol.AuthStatus) {
		buffer := buf.NewBuffer(protocol.AuthResponseLen)
		defer buffer.Release()

		log.Debug("Write auth response body, status:[%d]", status)
		response := protocol.NewAuthResponse(status)
		if err := response.WriteTo(buffer); err != nil {
			log.Warn("write to auth response error", err)
			return
		}
		if err := conn.ReadFrom(buffer); err != nil {
			log.Warn("Write to auth response error", err)
			return
		}
	}
)

func AuthToken(conn net.Conn) (*protocol.AuthHeader, *protocol.SessionAuth, error) {
	header, body, err := decodeAuthMsg(conn)
	if err != nil {
		return nil, nil, err
	}

	sessionAuth := body.(*protocol.SessionAuth)
	if err := VerificationToken(sessionAuth.Token); err != nil {
		complete(conn, protocol.AuthFail)
		return nil, nil, errors.NewError("Auth fail", err)
	}
	// success
	complete(conn, protocol.AuthSuccess)
	return header, sessionAuth, nil
}

func AuthKey(key string, conn net.Conn) (*protocol.AuthHeader, *protocol.ClientAuth, error) {
	header, body, err := decodeAuthMsg(conn)
	if err != nil {
		return nil, nil, err
	}

	clientAuth := body.(*protocol.ClientAuth)
	if common.IsNotEmpty(key) && clientAuth.Key != key {
		complete(conn, protocol.AuthFail)
		return nil, nil, errors.NewError("Auth fail, client key not match", nil)
	}

	// success
	complete(conn, protocol.AuthSuccess)
	return header, clientAuth, nil
}

// Authentication
func decodeAuthMsg(conn net.Conn) (*protocol.AuthHeader, proto.Message, error) {
	buffer := buf.NewBuffer(256)
	defer buffer.Release()

	if err := conn.SetDeadline(time.Now().Add(Timout)); err != nil {
		return nil, nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	if n, err := buffer.Write(conn, protocol.AuthHeaderLen); err != nil || n != protocol.AuthHeaderLen {
		return nil, nil, errors.NewError(
			fmt.Sprintf("Read bytes from connect too short, should minnum read %d bytes, actual reading %d bytes",
				protocol.AuthHeaderLen, n), err)
	}

	header := &protocol.AuthHeader{}
	if err := protocol.DecodeAuthHeader(buffer, header); err != nil {
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
	return header, body, nil
}
