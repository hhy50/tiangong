package session

import (
	"context"

	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ManagerCompName            = "SessionManager"
	sessions        []*Session = make([]*Session, 128)
)

type SessionManager struct {
}

func init() {
	component.Register(ManagerCompName, func(ctx context.Context) (component.Component, error) {
		return &SessionManager{}, nil
	})
}

func (s SessionManager) Start() error {
	return nil
}

func (s *SessionManager) AddSession(session *Session) error {
	sessions = append(sessions, session)
	log.Info("New session connected. token=%s, subHost=%s", session.Token, session.SubHost)
	return nil
}
