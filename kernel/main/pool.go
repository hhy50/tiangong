package main

import (
	"context"
	"math"
	"sync/atomic"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/kernel/tool"
)

const maxConnect = 1

var (
	// free list
	freeList = &tool.LinkedList{}
	// not active
	noActiveList = &tool.LinkedList{}

	//
	notifyGroup = []func(){}
	nLock       = lock.NewLock()

	incrementer        = common.Incrementer{Range: common.Range{0, math.MaxUint32}}
	reourceCount int32 = 0
)

func init() {
	net.DefaultConnTimeout = 10 * time.Second
}

func GetResource() (*Resource, error) {
	if res := freeList.Pop(); res != nil {
		return res.(*Resource), nil
	}

	for reourceCount < maxConnect {
		if !atomic.CompareAndSwapInt32(&reourceCount, reourceCount, reourceCount+1) {
			continue
		}
		resource := &Resource{
			TcpClient: net.NewTcpClient(Config.ServerHost, Config.ServerPort, context.Background()),
			num:       int(reourceCount),
			incre:     incrementer,
		}
		if err := resource.Connect(ConnSuccess); err != nil {
			noActiveList.Put(resource)
			return nil, errors.NewError("Connect to target server error. ", err)
		}
		return resource, nil
	}
	if freeList.Empty() {
		return nil, errors.NewError("There are no more resources to obtain.", nil)
	}
	return nil, errors.NewError("Connect to target server error. ", nil)
}

func GetResourceWithTimeout(timeout time.Duration) (r *Resource, err error) {
	if r, err = GetResource(); r != nil {
		return
	}

	now := time.Now()
	deadline := now.Add(timeout)

	for deadline.After(now) {
		ctx := JoinNotifyGooup(deadline.Sub(now))
		<-ctx.Done()
		if r, err = GetResource(); r != nil {
			return
		}
		now = time.Now()
	}
	return
}

func PutResource(res *Resource) {
	if res != nil {
		freeList.Put(res)
		// notify first
		nLock.Lock()
		defer nLock.Unlock()

		for len(notifyGroup) > 0 && !freeList.Empty() {
			notifyGroup[0]()
			notifyGroup = notifyGroup[1:]
		}
	}
}

func JoinNotifyGooup(timeout time.Duration) context.Context {
	nLock.Lock()
	defer nLock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	notifyGroup = append(notifyGroup, cancel)
	return ctx
}
