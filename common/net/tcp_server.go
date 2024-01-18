package net

import (
	"context"
	"fmt"
	"net"
	"tiangong/common"
	"tiangong/common/errors"
	"tiangong/common/log"
	"time"
)

var (
	logPrefix     = "[TCP]"
	AcceptTimeout = 5 * time.Second
)

type TcpServer interface {
	ListenTCP(handler ConnHandlerFunc) error
}

type tcpServerImpl struct {
	Host string
	Port Port

	listener net.Listener
	ctx      context.Context
	status   int
}

func (s *tcpServerImpl) ListenTCP(handler ConnHandlerFunc) error {
	if common.IsEmpty(s.Host) {
		s.Host = Local.String()
	}
	if s.Port == 0 {
		return errors.NewError("Server port not be null", nil)
	}
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	log.Info("%s Listen Host: %s, port: %d", logPrefix, s.Host, s.Port)
	go listenConnect(s, handler)
	return nil
}

func (s *tcpServerImpl) Stop() {
	if s.listener != nil {
		_ = s.listener.Close()
		s.listener = nil
	}
}

func listenConnect(s *tcpServerImpl, connHandler ConnHandlerFunc) {
	defer s.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if listener, ok := s.listener.(*net.TCPListener); ok {
				_ = listener.SetDeadline(time.Now().Add(AcceptTimeout))
			}
			conn, err := s.listener.Accept()
			if err != nil {
				continue
			}
			if err = connHandler(s.ctx, ConnWrap{conn}); err != nil {
				_ = conn.Close()
			}
		}
	}
}

func NewTcpServer(host string, port int, ctx context.Context) TcpServer {
	return &tcpServerImpl{
		Host: host,
		Port: Port(port),
		ctx:  ctx,
	}
}
