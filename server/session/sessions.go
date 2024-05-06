package session

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/lock"

	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ManagerName = "SessionManager"
)

type Manager struct {
	sessions []*Session

	lock lock.Lock
}

func init() {
	component.Register(ManagerName, func(ctx context.Context) (component.Component, error) {
		return &Manager{
			sessions: []*Session{},
			lock:     lock.NewLock(),
		}, nil
	})
}

func (sm Manager) Start() error {
	return nil
}

func (sm *Manager) AddSession(new *Session) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sm.sessions = append(sm.sessions, new)
	log.Info("New session connected. token=%s, subHost=%s", new.Token, new.SubHost)
	return nil
}

func (sm *Manager) Remove(session *Session) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	for i := range sm.sessions {
		if sm.sessions[i] == session {
			tail := len(sm.sessions) - 1
			sm.sessions[i], sm.sessions[tail] = sm.sessions[tail], sm.sessions[i]
			sm.sessions = sm.sessions[:tail]
			break
		}
	}
	log.Info("The session closed. token=%s, subHost=%s", session.Token, session.SubHost)
	return nil
}
