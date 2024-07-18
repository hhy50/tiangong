package common

import "github.com/haiyanghan/tiangong/common/log"

type FuncRunnable = func()

func SafeRun(runner FuncRunnable) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("goroutine panic, [%+v]", nil, err)
		}
	}()
	runner()
}

func SafeCall(runner FuncRunnable) {
	SafeRun(runner)
}
