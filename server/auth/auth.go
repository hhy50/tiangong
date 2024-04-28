package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
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

func AuthToken(conn net.Conn) (*protocol.AuthPacketHeader, *protocol.SessionAuthBody, error) {
	header, body, err := decodeAuthMsg(conn)
	if err != nil {
		return nil, nil, err
	}

	sessionAuthBody := body.(*protocol.SessionAuthBody)
	if err := VerificationToken(sessionAuthBody.Token); err != nil {
		complete(conn, protocol.AuthFail)
		return nil, nil, errors.NewError("Auth fail", err)
	}

	// success
	complete(conn, protocol.AuthSuccess)
	return header, sessionAuthBody, nil
}

func AuthKey(key string, conn net.Conn) (*protocol.AuthPacketHeader, *protocol.ClientAuthBody, error) {
	clientAuth, body, err := decodeAuthMsg(conn)
	if err != nil {
		return nil, nil, err
	}

	clientAuthBody := body.(*protocol.ClientAuthBody)
	if common.IsNotEmpty(key) && clientAuthBody.Key != key {
		complete(conn, protocol.AuthFail)
		return nil, nil, errors.NewError("Auth fail, client key not match", nil)
	}

	// success
	complete(conn, protocol.AuthSuccess)
	return clientAuth, clientAuthBody, nil
}

// decodeAuthMsg
func decodeAuthMsg(conn net.Conn) (*protocol.AuthPacketHeader, interface{}, error) {
	buffer := buf.NewBuffer(256)
	defer buffer.Release()

	if err := conn.SetDeadline(time.Now().Add(Timout)); err != nil {
		return nil, nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	if n, err := buffer.Write(conn, protocol.PacketHeaderLen); err != nil || n != protocol.PacketHeaderLen {
		return nil, nil, errors.NewError(
			fmt.Sprintf("Read bytes from connect too short, should minnum read %d bytes, actual reading %d bytes",
				protocol.PacketHeaderLen, n), err)
	}

	header := &protocol.AuthPacketHeader{}
	if err := header.ReadFrom(buffer); err != nil {
		return nil, nil, err
	}

	size := int(header.Len)
	if n, err := buffer.Write(conn, size); err != nil || n != size {
		return nil, nil, errors.NewError(
			fmt.Sprintf("Read auth body from connect too short, should minnum read %d bytes, actual reading %d bytes", size, n), err)
	}

	var body interface{}
	switch header.AuthType() {
	case protocol.AuthClient:
		body = &protocol.ClientAuthBody{}
	case protocol.AuthSession:
		body = &protocol.SessionAuthBody{}
	}

	bytes, _ := buf.ReadAll(buffer)
	if err := json.Unmarshal(bytes, body); err != nil {
		return nil, nil, err
	}
	return header, body, nil
}
