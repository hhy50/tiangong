package server

import (
	"os"
	"path/filepath"
	"tiangong/common/errors"
	"tiangong/common/lock"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/server/admin"
	"tiangong/server/client"
	"tiangong/server/conf"
)

type ServerStatus int8

type Server interface {
	Start() error
	Stop()
}

type tgServer struct {
	Admin   admin.AdminServer
	Clients map[string]*client.Client
	Status  ServerStatus
	Lock    lock.Lock

	TcpSrv net.TcpServer
}

const (
	INIT ServerStatus = iota
	RUNNING
	STOPED
)

var (
	connHandler = client.ConnHandler

	getRedomPasswd = func() string {
		exec, err := os.Executable()
		if err != nil {
			return ""
		}
		return filepath.Dir(exec)
	}
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
	log.Warn("TianGong Server end...")
}

func NewServer(input string) (Server, error) {
	config, err := conf.LoadConfig(input)
	if err != nil {
		return nil, err
	}

	admin := admin.AdminServer{
		HttpPort: config.HttpPort,
		UserName: config.UserName,
		Password: config.Passwd,
	}

	tcpSrv := net.TcpServer{
		Host: config.Host,
		Port: config.SrvPort,
	}

	svr := &tgServer{
		Admin:   admin,
		Clients: make(map[string]*client.Client),
		Status:  INIT,
		TcpSrv:  tcpSrv,
		Lock:    lock.NewLock(),
	}
	return svr, nil
}
