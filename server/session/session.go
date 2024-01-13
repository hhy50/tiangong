package session

import (
	"tiangong/common/buf"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

type Session struct {
	Token   string
	SubHost net.IpAddress

	buffer buf.Buffer
	conn   net.Conn
}

//
// +----+------+----------+
// | 	PacketHeader      |
// +----+------+----------+
// |Len | Rid  |  Protol  |
// +----+------+----------+
// | 2  |  4   | 	1     |
// +----+------+----------+
func (s *Session) Work() {
	defer s.Close()
	if err := s.conn.SetDeadline(time.Time{}); err != nil {
		return
	}
	for {
		header := make([]byte, protocol.PacketHeaderLen)
		if _, err := s.conn.Read(header); err != nil {
			log.Error("Read error from session, ", err)
			return
		}
		packetHeader := protocol.PacketHeader{}
		if err := packetHeader.Unmarshal(header); err != nil {
			s.Close()
		}
		packetLen := packetHeader.Len
		if n, err := s.buffer.Write(s.conn, int(packetLen)); err != nil || n != int(packetLen) {
			// discard
			_ = s.buffer.Clear()
			continue
		}

		// TODO
	}
}

func (s *Session) Close() {
	s.buffer.Release()
	_ = s.conn.Close()
	log.Warn("Session Closed, token: %s", s.Token)
}
