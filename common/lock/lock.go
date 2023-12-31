package lock

import (
	"sync"
	"time"
)

var (
	GlobalLock = Mutex{}
)

type Lock interface {
	Lock()
	TryLock(time.Duration) bool
	Unlock()
}

type Mutex struct {
	mutex sync.Mutex
}

func (m *Mutex) Lock() {
	m.mutex.Lock()
}

func (m *Mutex) Unlock() {
	m.mutex.Unlock()
}

func (m *Mutex) TryLock(t time.Duration) bool {
	timer := time.After(t)
	for {
		select {
		case <-timer:
			return false
		default:
			if m.mutex.TryLock() {
				return true
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func NewLock() Lock {
	return &Mutex{}
}
