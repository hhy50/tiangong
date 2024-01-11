package session

import (
	"tiangong/common/buf"
	"tiangong/common/net"
)

var sessions []*Session = make([]*Session, 128)

func AddSession(session *Session) error {
	sessions = append(sessions, session)
	return nil
}

func NewSession(subHost net.IpAddress, token string, conn net.Conn) Session {
	return Session{
		SubHost: subHost,
		Token:   token,

		buffer: buf.NewRingBuffer(),
		conn:   conn,
	}
}
