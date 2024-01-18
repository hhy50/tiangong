package net

import (
	"context"
	"fmt"
	"net"
	"tiangong/common"
	"tiangong/common/log"
)

var logPrefix = "[TCP]"

type TcpServer struct {
	Host string
	Port int

	ctx    common.Context
	cancel func()
}

func (s *TcpServer) Listen(handler ConnHandlerFunc) error {
	ctx, cancel := context.WithCancel(common.EmptyCtx)
	s.ctx, s.cancel = common.Wrap(ctx), cancel

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	go listenConnect(listener, handler, s.ctx)

	s.ctx.Add(common.ListenerCtxKey, listener)
	log.Info("%s Listen Host: %s, port: %d", logPrefix, s.Host, s.Port)
	return nil
}

func (s *TcpServer) Stop() {
	s.cancel()
	if listener, ok := s.ctx.Value(common.ListenerCtxKey).(net.Listener); ok {
		log.Warn("%s listener stopping...", logPrefix)
		_ = listener.Close()
	}
}

func listenConnect(listener net.Listener, connHandler ConnHandlerFunc, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			if err = connHandler(ConnWrap{conn}); err != nil {
				_ = conn.Close()
			}
		}
	}
}
