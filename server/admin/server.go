package admin

import "tiangong/common/net"

type AdminServer struct {
	HttpPort int
	UserName string
	Password string

	tcpServer net.TcpServer
}

func (s *AdminServer) Start() error {
	return nil
}

func (s *AdminServer) Stop() {

}
