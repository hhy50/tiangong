package common

import (
	"tiangong/common/log"
	"time"
)

type Retry func() error

func (r Retry) Run(interval, timeout time.Duration) {
	always := timeout < 0
	to := time.Now().Add(timeout)
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		if err := r(); err != nil {
			log.Error("retry execute error", err)
			if !always && time.Now().After(to) {
				ticker.Stop()
				break
			}
		} else {
			ticker.Stop()
			break
		}
	}
}
