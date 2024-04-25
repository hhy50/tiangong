//go:build debug

package common

import (
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common/log"
)

func init() {
	go TimerFunc(func() {
		buf := make([]byte, 4096)
		length := runtime.Stack(buf, true)

		monitor := "The number of active goroutine in the program: %d, detail: \n %s"
		log.Debug(monitor, runtime.NumGoroutine(), String(buf[:length]))
	}).Run(60 * time.Second)
}
