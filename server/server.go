package server

import (
	_ "github.com/haiyanghan/tiangong/server/admin"
	_ "github.com/haiyanghan/tiangong/server/client"
	_ "github.com/haiyanghan/tiangong/server/session"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	Running = 1
)

type Status int8
type Runnable = common.Runnable

type Server interface {
	Start() error
	Stop()
}

type tgServer struct {
	Lock lock.Lock
	ctx  context.Context

	status int
}

func (s *tgServer) Start() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.ctx.AddValue("Instance", s)
	for name, _ := range component.GetComponents() {
		value := s.ctx.Value(name)
		if value == nil {
			continue
		}
		comp := value.(component.Component)
		if err := comp.Start(); err != nil {
			return err
		}
	}
	s.status = Running
	return nil
}

func (s *tgServer) Stop() {
	if s.status != Running {
		return
	}
	s.ctx.Cancel()
	log.Warn("TianGong Server end...")
}

func NewServer() (Server, error) {
	ctx := context.Empty()
	for name, compCtreator := range component.GetComponents() {
		comp, err := compCtreator(ctx)
		if err != nil {
			return nil, err
		}
		ctx.AddValue(name, comp)
	}

	return &tgServer{
		Lock: lock.NewLock(),
		ctx:  ctx,
	}, nil
}
