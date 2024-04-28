package session

import (
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/server/client"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Session struct {
	Token string

	ctx    context.Context
	buffer buf.Buffer
	conn   net.Conn
	bridge Bridge
}

func (s *Session) Work() {
	defer s.Close()
	for {
		select {
		case <-s.ctx.Done():
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

// HandlePacket
func (s *Session) HandlePacket() error {
	if _, err := s.buffer.Write(s.conn, protocol.PacketHeaderLen); err != nil {
		return err
	}
	header := protocol.DataPacket{}
	if err := header.ReadFrom(s.buffer); err != nil {
		return err
	}
	if header.Len == 0 {
		return nil
	}
	log.Debug("Receive packet header, rid:%d, len:%d", header.Rid, header.Len)
	if n, err := s.buffer.Write(s.conn, int(header.Len)); err != nil || n != int(header.Len) {
		// discard
		_ = s.buffer.Clear()
		return nil
	}
	if err := s.bridge.Transport(&header, s.buffer); err != nil {
		return err
	}
	return nil
}

func NewSession(token string, conn net.Conn, ctx context.Context, dstClient *client.Client) Session {
	return Session{
		Token:  token,
		ctx:    ctx,
		conn:   conn,
		bridge: &WirelessBridging{dstClient},
		buffer: buf.NewRingBuffer(),
	}
}
