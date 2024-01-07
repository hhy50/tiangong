package auth

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"net"
	"tiangong/common"
	"tiangong/common/errors"
	tgNet "tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"tiangong/server"
	"tiangong/server/client"
	"tiangong/server/internal"
	"time"
)

const (
	Client_Auth  protocol.AuthType = 1
	Session_Auth protocol.AuthType = 2
)

var (
	TimeOut = 15 * time.Second
)

type ListenFunc func()

func (l ListenFunc) Run() { l() }

func Authentication(conn net.Conn) (server.Runnable, error) {
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
		return nil, err
	}

	// 验证服务端key的有效性
	switch body.(type) {
	case *protocol.ClientAuth:
		clientAuth := body.(*protocol.ClientAuth)
		if clientAuth.Key != server.Key {
			return nil, errors.NewError("Auth fail, client key not match", nil)
		}
		c := buildClient(conn, header, clientAuth)
		_ = client.AddClient(&c)
		return ListenFunc(func() {

		}), nil
	case *protocol.SessionAuth:
		sessionAuth := body.(*protocol.SessionAuth)
		if err = Verification(sessionAuth.Token); err != nil {
			return nil, errors.NewError("Auth fail", err)
		}
		buildSession(conn, header, sessionAuth)
	}
	return nil, nil
}

func buildClient(conn net.Conn, header *protocol.AuthHeader, auth *protocol.ClientAuth) client.Client {
	getInternalIpFromReq := func() tgNet.IpAddress {
		if len(auth.Internal) == 4 {
			return auth.Internal[0:4]
		}
		return nil
	}

	internalIp := getInternalIpFromReq()
	if internalIp != nil {
		internalIp = internal.GeneraInternalIp()
	}
	clientName := auth.Name
	if common.IsEmpty(clientName) {
		uid, _ := uuid.NewUUID()
		clientName = uid.String()
	}
	return client.NewClient(clientName, internalIp, header, conn)
}

func buildSession(net.Conn, *protocol.AuthHeader, *protocol.SessionAuth) {

}
