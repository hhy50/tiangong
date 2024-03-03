package server

import (
	"context"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/admin"
	"github.com/haiyanghan/tiangong/server/client"
)

var (
	ServerCnf Config
	Running   = 1
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
	Lock    lock.Lock
	TcpSrv  net.TcpServer
	Ctx     context.Context

	status int
}

func (s *tgServer) Start() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if err := s.TcpSrv.ListenTCP(connHandler); err != nil {
		return err
	}
	s.status = Running
	return nil
}

func (s *tgServer) Stop() {
	if s.status != Running {
		return
	}
	cancel := s.Ctx.Value(common.CancelFuncKey).(context.CancelFunc)
	cancel()
	s.status = 0
	log.Warn("TianGong Server end...")
}

func NewServer(input string) (Server, error) {
	if err := conf.LoadConfig(input, &ServerCnf, defaultValue); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, common.CancelFuncKey, cancel)

	adm := admin.AdminServer{
		HttpPort: ServerCnf.HttpPort,
		UserName: ServerCnf.UserName,
		Password: ServerCnf.Passwd,
		Ctx:      ctx,
	}

	server := &tgServer{
		Admin:   adm,
		Clients: make(map[string]*client.Client),
		TcpSrv:  net.NewTcpServer(ServerCnf.Host, ServerCnf.SrvPort, ctx),
		Lock:    lock.NewLock(),
		Ctx:     ctx,
	}
	return server, nil
}
