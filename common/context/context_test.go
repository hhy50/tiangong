package context_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/haiyanghan/tiangong/common/context"
)

func TestGetValue(t *testing.T) {
	ctx0 := context.Empty()
	ctx0.AddValue("a", "a")
	ctx0 = context.WithParent(ctx0)
	ctx0.AddValue("b", "b")
	ctx0 = context.WithParent(ctx0)
	ctx0.AddValue("c", "c")
	ctx0 = context.WithParent(ctx0)
	ctx0.AddValue("d", "d")
	ctx0 = context.WithParent(ctx0)
	ctx0.AddValue("e", "e")
	ctx0 = context.WithParent(ctx0)

	if a := ctx0.Value("a"); a == nil || a != "a" {
		t.Error("get value from context error, not found key 'a' ")
		return
	}
	if a := ctx0.Value("b"); a == nil || a != "b" {
		t.Error("get value from context error, not found key 'b' ")
		return
	}
	if a := ctx0.Value("c"); a == nil || a != "c" {
		t.Error("get value from context error, not found key 'c' ")
		return
	}
	if a := ctx0.Value("d"); a == nil || a != "d" {
		t.Error("get value from context error, not found key 'd' ")
		return
	}
	if a := ctx0.Value("e"); a == nil || a != "e" {
		t.Error("get value from context error, not found key 'e' ")
		return
	}
}

func TestCancel(t *testing.T) {
	ctx0_cancel := atomic.Bool{}
	ctx0 := context.Empty()

	go func() {
		<-ctx0.Done()
		ctx0_cancel.Store(true)
	}()

	ctx1_cancel := atomic.Bool{}
	ctx1 := context.Empty()

	go func() {
		<-ctx1.Done()
		ctx1_cancel.Store(true)
	}()

	ctx2_cancel := atomic.Bool{}
	ctx2 := context.Empty()

	go func() {
		<-ctx2.Done()
		ctx2_cancel.Store(true)
	}()

	ctx3_cancel := atomic.Bool{}
	ctx3 := context.Empty()

	go func() {
		<-ctx3.Done()
		ctx3_cancel.Store(true)
	}()
	ctx3.Cancel()
	time.Sleep(time.Second)

	if !ctx3_cancel.Load() {
		t.Error("ctx3 is not cancel")
		return
	}

	if ctx1_cancel.Load() || ctx2_cancel.Load() {
		t.Error("ctx1 | ctx2 unexpected shutdown")
		return
	}

	ctx1.Cancel()
	time.Sleep(time.Second)

	if !ctx1_cancel.Load() || !ctx2_cancel.Load() {
		t.Error("ctx1 | ctx2 is not cancel")
		return
	}
}
