package common

import "time"

type TimerFunc func()

type OnecTimerFunc func()

func (t TimerFunc) Run(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	
	for {
		<-ticker.C
		t()
	}
}

func (t OnecTimerFunc) Run(d time.Duration) {
	after := time.After(d)
	<-after
	t()
}
