package common

import "time"

type TimerFunc func()

type OnceTimerFunc func()

func (t TimerFunc) Run(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		<-ticker.C
		t()
	}
}

func (t OnceTimerFunc) Run(d time.Duration) {
	after := time.After(d)
	<-after
	t()
}
