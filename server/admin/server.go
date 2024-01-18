package admin

import (
	"context"
	"tiangong/common/net"
)

type AdminServer struct {
	HttpPort int
	UserName string
	Password string
	Ctx      context.Context

	tcpServer net.TcpServer
}

func (s *AdminServer) Start() error {
	return nil
}

func (s *AdminServer) Stop() {

}
