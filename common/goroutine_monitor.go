//go:build debug

package common

import (
	"github.com/haiyanghan/tiangong/common/log"
	"runtime"
	"time"
)

func init() {
	go TimerFunc(func() {
		buf := make([]byte, 4096)
		length := runtime.Stack(buf, true)

		monitor := "The number of active goroutine in the program: %d, detail: \n %s"
		log.Debug(monitor, runtime.NumGoroutine(), string(buf[:length]))
	}).Run(100 * time.Second)
}
