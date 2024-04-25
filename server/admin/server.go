package admin

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"
)

type AdminServer struct {
	HttpPort int
	UserName string
	Password string

	ctx       context.Context
	tcpServer net.TcpServer
}

func (s *AdminServer) Start() error {
	return nil
}

func (s *AdminServer) Stop() {

}
