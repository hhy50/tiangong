package server

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"tiangong/server/conf"
)

var getRedomPasswd = func() string {
	exec, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(exec)
}

type Server struct {
	conf *conf.ServerConfig
}

func (s *Server) Start() error {
	conf := s.conf
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.TcpPort))
	if err != nil {
		return err
	}
	fmt.Printf("listener host:%s, listener tcp port:%d, listener http port:%d \n", conf.Host, conf.TcpPort, conf.HttpPort)
	go func() {
		for {
			_, err := listener.Accept()
			if err != nil {

			}
		}
	}()
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func NewServer(cp string) (*Server, error) {
	conf, err := conf.LoadConfigWithPath(cp)
	if err != nil {
		return nil, err
	}
	server := &Server{
		conf: conf,
	}
	return server, nil
}
