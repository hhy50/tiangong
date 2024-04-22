package server

import (
	"context"

	_ "github.com/haiyanghan/tiangong/server/admin"
	_ "github.com/haiyanghan/tiangong/server/client"
	_ "github.com/haiyanghan/tiangong/server/session"

	"github.com/haiyanghan/tiangong/common"
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
	Lock  lock.Lock
	Comps map[string]component.Component
	Ctx   context.Context

	status int
}

func (s *tgServer) Start() error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	for _, comp := range s.Comps {
		comp.Start()
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

func NewServer() (Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, common.CancelFuncKey, cancel)

	comps := map[string]component.Component{}
	for name, compCtreator := range component.GetComponents() {
		comp, err := compCtreator(ctx)
		if err != nil {
			return nil, err
		}
		comps[name] = comp
		ctx = context.WithValue(ctx, name, comp)
	}
	return &tgServer{
		Comps: comps,
		Lock:  lock.NewLock(),
		Ctx:   ctx,
	}, nil
}
