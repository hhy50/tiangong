package session

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"

	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ManagerName = "SessionManager"
)

type Manager struct {
	sessions map[string]*Session
}

func init() {
	component.Register(ManagerName, func(ctx context.Context) (component.Component, error) {
		return &Manager{
			sessions: map[string]*Session{},
		}, nil
	})
}

func (s Manager) Start() error {
	return nil
}

func (s *Manager) AddSession(subhost net.IpAddress, new *Session) error {
	s.sessions[new.Token] = new
	log.Info("New session connected. token=%s, subHost=%s", new.Token, subhost)
	return nil
}

func (s *Manager) Remove(session *Session) error {
	delete(s.sessions, session.Token)
	log.Info("The session closed. token=%s, subHost=%s", session.Token)
	return nil
}
