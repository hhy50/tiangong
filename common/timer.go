package common

import "time"

type TimerFunc func()

type OnceTimerFunc func()

func (t TimerFunc) Run(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	defer Recover()

	for {
		<-ticker.C
		t()
	}
}

func (t OnceTimerFunc) Run(d time.Duration) {
	defer Recover()

	after := time.After(d)
	<-after
	t()
}
