package router

import (
	"tiangong/common/errors"
	"tiangong/common/lock"
	"tiangong/server/gateway"
)

var (
	RouterTable = make(map[string]*gateway.Destination)
	RTLock      = lock.NewLock()
)

func RegisterRouter(host string, dest *gateway.Destination) error {
	RTLock.Lock()
	defer RTLock.Unlock()

	if _, f := RouterTable[host]; f {
		return errors.NewError("Existing routing information: "+host, nil)
	}
	RouterTable[host] = dest
	return nil
}
