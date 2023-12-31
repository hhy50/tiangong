package net

import (
	"fmt"
	"net"
	"tiangong/common/context"
	"tiangong/common/log"
)

type TcpServer struct {
	Host string
	Port int

	ctx    context.Context
	cancel func()
}

func (s *TcpServer) Listen(handler ConnHandler) error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	go ListenConn(listener, handler, s.ctx)
	log.Info("Listen Host: %s, port: %d", s.Host, s.Port)
	return nil
}

func (s *TcpServer) Stop() {
	s.cancel()
}

func ListenConn(listener net.Listener, connHander ConnHandler, ctx context.Context) {
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			connHander(conn)
		}
	}

}
