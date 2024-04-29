package auth

import (
	"encoding/json"
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
		buffer := buf.NewBuffer(protocol.PacketHeaderLen)
		defer buffer.Release()

		packet := protocol.NewAuthResponsePacket(status)
		if err := protocol.EncodePacket(buffer, packet); err != nil {
			log.Warn("write to auth response error, %+v", err)
			return
		}
		if err := conn.ReadFrom(buffer); err != nil {
			log.Warn("Write to auth response error, %+v", err)
			return
		}
	}
)

func AuthToken(conn net.Conn) (*protocol.PacketHeader, *protocol.SessionAuthBody, error) {
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

func AuthKey(key string, conn net.Conn) (*protocol.PacketHeader, *protocol.ClientAuthBody, error) {
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
func decodeAuthMsg(conn net.Conn) (*protocol.PacketHeader, interface{}, error) {
	buffer := buf.NewBuffer(4096)
	defer buffer.Release()

	if err := conn.SetDeadline(time.Now().Add(Timout)); err != nil {
		return nil, nil, errors.NewError("Auth fail, SetDeadline error", err)
	}

	packet, err := protocol.DecodePacket(buffer, conn)
	if err != nil {
		return nil, nil, err
	}
	var body interface{}
	switch packet.AuthType() {
	case protocol.AuthSession:
		body = &protocol.SessionAuthBody{}
	case protocol.AuthClient:
		body = &protocol.ClientAuthBody{}
	}
	if err := json.Unmarshal(packet.Body, body); err != nil {
		return nil, nil, err
	}
	return &packet.Header, body, nil
}
