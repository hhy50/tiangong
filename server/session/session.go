package session

import (
	"context"
	"runtime"
	"tiangong/common/buf"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

type Session struct {
	Token   string
	SubHost net.IpAddress
	Ctx     context.Context

	buffer buf.Buffer
	conn   net.Conn
	bridge Bridge
}

// +----+------+----------+
// | 	PacketHeader      |
// +----+------+----------+
// |Len | Rid  |  Protol  |
// +----+------+----------+
// | 2  |  4   | 	1     |
// +----+------+----------+
func (s *Session) Work() {
	defer s.Close()
	select {
	case <-s.Ctx.Done():
		runtime.Goexit()
	default:
		if err := s.conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
			log.Error("SetDeadline error", err)
			return
		}
		if err := s.HandlePacket(); err != nil {
			log.Error("HandlePacket error, ", err)
			return
		}
	}
}

func (s *Session) Close() {
	s.buffer.Release()
	_ = s.conn.Close()
	log.Warn("Session Closed, token: %s", s.Token)
}

func (s *Session) HandlePacket() error {
	if _, err := s.buffer.Write(s.conn, protocol.PacketHeaderLen); err != nil {
		return err
	}
	header := protocol.PacketHeader{}
	if err := header.ReadFrom(s.buffer); err != nil {
		s.Close()
	}
	log.Debug("Receive packet header, protocol:%s, rid:%d, len:%d", protocol.Protocol(header.Protocol).String(), header.Rid, header.Len)
	if n, err := s.buffer.Write(s.conn, int(header.Len)); err != nil || n != int(header.Len) {
		// discard
		_ = s.buffer.Clear()
		log.Warn("Discard packet, len:%d, error:%+v", n, err)
		return nil
	}

	if err := s.bridge.Transport(header, s.buffer); err != nil {
		return err
	}
	return nil
}
