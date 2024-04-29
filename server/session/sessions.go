package session

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"

	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ManagerCompName = "SessionManager"
)

type Manager struct {
	sessions []*Session
}

func init() {
	component.Register(ManagerCompName, func(ctx context.Context) (component.Component, error) {
		return &Manager{
			sessions: make([]*Session, 128),
		}, nil
	})
}

func (s Manager) Start() error {
	return nil
}

func (s *Manager) AddSession(subhost net.IpAddress, session *Session) error {
	s.sessions = append(s.sessions, session)
	log.Info("New session connected. token=%s, subHost=%s", session.Token, subhost)
	return nil
}
