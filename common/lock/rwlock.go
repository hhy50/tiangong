package lock

import "sync"

// TODO

type Rwlock struct {
	sync.RWMutex
}
