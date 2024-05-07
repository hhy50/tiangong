package common

import (
	"time"

	"github.com/haiyanghan/tiangong/common/log"
)

type TimerFunc func()

type OnceTimerFunc func()

func (t TimerFunc) Run(d time.Duration) {
	ticker := time.NewTicker(d)
	defer func() {
		if err := recover(); err != nil {
			log.Error("goroutine panic, %+v", nil, err)
		}
		ticker.Stop()
	}()

	for {
		<-ticker.C
		t()
	}
}

func (t OnceTimerFunc) Run(d time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("goroutine panic, %+v", nil, err)
		}
	}()

	after := time.After(d)
	<-after
	t()
}
