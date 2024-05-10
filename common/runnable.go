package common

import "github.com/haiyanghan/tiangong/common/log"

type Runnable interface {
	Run()
}

type FuncRunable func()

func (fn FuncRunable) Run() {
	fn()
}

func SafeRun(runner Runnable) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("goroutine panic, [%+v]", nil, err)
		}
	}()
	runner.Run()
}

func SafeCall(runner FuncRunable) {
	SafeRun(runner)
}
