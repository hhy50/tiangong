package session

import (
	"context"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Session struct {
	Token   string
	SubHost net.IpAddress
	Ctx     context.Context

	buffer buf.Buffer
	conn   net.Conn
	bridge Bridge
}

// +----+-----+--------+----------+--------+
// |	    PacketHeader (20 byte)		   |
// +----+-----+--------+----------+--------+
// | Len | Rid | Protol | Reserved | Status |
// +----+-----+--------+----------+--------+
// | 2  |  4  |   1    |   12	  |   1    |
// +----+-----+--------+----------+--------+
func (s *Session) Work() {
	defer s.Close()
	for {
		select {
		case <-s.Ctx.Done():
			runtime.Goexit()
		default:
			if err := s.conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				log.Error("SetDeadline error", err)
				return
			}
			if err := s.HandlePacket(); err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				log.Warn("HandlePacket error, %v", err)
				return
			}
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
		return err
	}
	if header.Len == 0 {
		return nil
	}
	log.Debug("Receive packet header, protocol:%s, rid:%d, len:%d", protocol.ProtocolToStr(header.Protocol), header.Rid, header.Len)
	// if n, err := s.buffer.Write(s.conn, int(header.Len)); err != nil || n != int(header.Len) {
	// 	// discard
	// 	discard(s.conn, int(header.Len)-n)
	// 	_ = s.buffer.Clear()
	// 	return nil
	// }
	if err := s.bridge.Transport(header, s.buffer); err != nil {
		return err
	}
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

func discard(conn net.Conn, len int) {
	discard := buf.NewBuffer(len)
	discard.Write(conn, len)
	discard.Release()

	log.Warn("Discard packet, len:%d", len)
}
