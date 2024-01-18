package server

import (
	"tiangong/common"
	"tiangong/common/conf"
	"tiangong/common/lock"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/server/admin"
	"tiangong/server/client"
)

type Status int8
type Runnable = common.Runnable

type Server interface {
	Start() error
	Stop()
}

type tgServer struct {
	Cnf     Config
	Admin   admin.AdminServer
	Clients map[string]*client.Client
	Lock    lock.Lock

	TcpSrv net.TcpServer
}

func (s *tgServer) Start() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if err := s.TcpSrv.Listen(s.connHandler); err != nil {
		return err
	}

	return nil
}

func (s *tgServer) Stop() {
	s.Admin.Stop()
	s.TcpSrv.Stop()
	log.Warn("TianGong Server end...")
}

func (s *tgServer) AddMapping(src string, sp int, dest string, dp int) error {

	return nil
}

func NewServer(input string) (Server, error) {
	config := Config{}
	if err := conf.LoadConfig(input, &config, defaultValue); err != nil {
		return nil, err
	}

	adm := admin.AdminServer{
		HttpPort: config.HttpPort,
		UserName: config.UserName,
		Password: config.Passwd,
	}

	tcpSrv := net.TcpServer{
		Host: config.Host,
		Port: config.SrvPort,
	}

	svr := &tgServer{
		Cnf:     config,
		Admin:   adm,
		Clients: make(map[string]*client.Client),
		TcpSrv:  tcpSrv,
		Lock:    lock.NewLock(),
	}
	return svr, nil
}
