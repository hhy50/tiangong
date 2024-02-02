package main

import (
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/log"
)

func init() {
	go common.TimerFunc(func() {
		log.Debug("Free linked count:%d, No active linked count:%d", freeList.Len(), noActiveList.Len())
	}).Run(time.Second)
}
