package common

import (
	"time"
)

type Retry func() error

func (r Retry) Run(interval, timeout time.Duration) {
	always := timeout < 0
	to := time.Now().Add(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := r(); err != nil {
			if !always && time.Now().After(to) {
				break
			}
		} else {
			break
		}
		<-ticker.C
	}
}
