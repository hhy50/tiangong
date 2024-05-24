package session

import (
	"runtime"
	"sync"
	"time"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/server/client"

	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Session struct {
	Token   string
	SubHost string
	ctx     context.Context
	bridge  Bridge
	once    sync.Once
}

func (s *Session) Work() {
	buffer := buf.NewRingBuffer()
	conn := s.ctx.Value(net.ConnValName).(net.Conn)

	defer func() {
		_ = conn.Close()
		buffer.Release()
		s.Close()
	}()

	for {
		select {
		case <-s.ctx.Done():
			runtime.Goexit()
		default:
			if packet, err := protocol.DecodePacket(buffer, conn, time.Minute); err != nil {
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
	s.once.Do(func() {
		// remove from manager
		sm := s.ctx.Value(ManagerName).(*Manager)
		sm.Remove(s)

		s.ctx.Cancel()
		log.Warn("Session Closed, token: %s", s.Token)
	})
}

// HandlePacket
func (s *Session) handlePacket(packet *protocol.Packet) error {
	log.Debug("Receive packet, rid:%d, len:%d", packet.Header.Rid, packet.Header.Len)
	if err := s.bridge.Transport(packet); err != nil {
		return err
	}
	return nil
}

func NewSession(ctx context.Context, token, subHot string, dstClient *client.Client) *Session {
	ctx = context.WithParent(ctx)
	var bridge Bridge
	if dstClient.Name == "Default" {
		bridge = &DirctClientBridging{
			ctx: ctx,
		}
	} else {
		bridge = &WirelessBridging{dstClient}
	}

	return &Session{
		Token:   token,
		SubHost: subHot,
		ctx:     ctx,
		bridge:  bridge,
	}
}
