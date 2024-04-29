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
			if packet, err := protocol.DecodePacket(s.buffer, s.conn); err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				log.Error("DecodePacket error, close the session", err)
				return
			} else if err = s.handlePacket(packet); err != nil {
				log.Error("handlePacket error", err)
			}
		}
	}
}

func (s *Session) Close() {
	_ = s.conn.Close()
	s.buffer.Release()
	log.Warn("Session Closed, token: %s", s.Token)
}

// HandlePacket
func (s *Session) handlePacket(packet *protocol.Packet) error {
	log.Debug("Receive packet, rid:%d, len:%d", packet.Header.Rid, packet.Header.Len)
	if err := s.bridge.Transport(packet, s.buffer); err != nil {
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
