package common

import (
	"time"
)

type TimerFunc func()

type OnceTimerFunc func()

func (t TimerFunc) Run(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		<-ticker.C
		SafeCall(FuncRunable(t))
	}
}

func (t OnceTimerFunc) Run(d time.Duration) {
	after := time.After(d)
	<-after
	SafeCall(FuncRunable(t))
}
