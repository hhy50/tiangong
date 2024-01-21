package session

import (
	"context"
	"tiangong/common/buf"
	"tiangong/common/net"
	"tiangong/server/client"
)

var sessions []*Session = make([]*Session, 128)

func AddSession(session *Session) error {
	sessions = append(sessions, session)
	return nil
}

func NewSession(subHost net.IpAddress, token string, conn net.Conn, ctx context.Context) Session {
	return Session{
		SubHost: subHost,
		Token:   token,
		Ctx:     ctx,

		bridge: &WirelessBridging{client.Clients[subHost]},
		buffer: buf.NewRingBuffer(),
		conn:   conn,
	}
}
