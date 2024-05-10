package common_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/log"
)

var (
	o = struct {
		F1 func()
	}{
		F1: func() {
			panic("panic")
		},
	}
)

func TestSafeCall(t *testing.T) {
	log.InitLog()

	crush := atomic.Bool{}
	defer func() {
		if !crush.Load() {
			t.Error("processor crush")
		}
	}()

	go common.SafeCall(o.F1)
	time.Sleep(time.Second)
	crush.Store(true)
}
