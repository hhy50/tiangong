package server

import (
	"tiangong/common"
	"tiangong/common/errors"
	"tiangong/common/lock"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/server/admin"
	"tiangong/server/client"
	"tiangong/server/conf"
)

type Status int8
type Runnable = common.Runnable

type Server interface {
	Start() error
	Stop()
}

type tgServer struct {
	Admin   admin.AdminServer
	Clients map[string]*client.Client
	Status  Status
	Lock    lock.Lock

	TcpSrv net.TcpServer
}

const (
	INIT Status = iota
	RUNNING
	STOPED
)

var (
	Key         string
	connHandler = client.ConnHandler
)

func (s *tgServer) Start() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if s.Status != INIT {
		return errors.NewError("Duplicate invoke start() error", nil)
	}
	if err := s.TcpSrv.Listen(connHandler); err != nil {
		return err
	}

	s.Status = RUNNING
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
	config, err := conf.LoadConfig(input)
	if err != nil {
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
		Admin:   adm,
		Clients: make(map[string]*client.Client),
		Status:  INIT,
		TcpSrv:  tcpSrv,
		Lock:    lock.NewLock(),
	}
	return svr, nil
}
