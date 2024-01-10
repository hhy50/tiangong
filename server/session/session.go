package session

import (
	"google.golang.org/protobuf/proto"
	"tiangong/common/buf"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
)

type Session struct {
	Host    net.IpAddress
	SubHost net.IpAddress
	Token   string

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
	for {
		header := make([]byte, protocol.PacketHeaderLen)
		if _, err := s.conn.Read(header); err != nil {
			log.Error("Read error from session, reason: %+v", err)
			continue
		}
		packetHeader := protocol.PacketHeader{}
		if err := proto.Unmarshal(header, &packetHeader); err != nil {
			s.Close()
		}
		if n, err := s.buffer.Write(s.conn); err != nil && n > 0 {

		}
	}
}

func (s *Session) Close() {
	_ = s.buffer.Release()
	_ = s.conn.Close()
}
