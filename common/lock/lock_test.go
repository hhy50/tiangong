package lock_test

import (
	"testing"
	"tiangong/common/lock"
	"time"
)

func TestGLock(t *testing.T) {
	start := time.Now().Unix()
	c := make(chan int64)
	go func() {
		lock.GlobalLock.Lock()
		c <- time.Now().Unix()
		time.Sleep(3 * time.Second)
		defer func() {
			lock.GlobalLock.Unlock()
		}()
	}()
	<-c
	lock.GlobalLock.Lock()
	defer lock.GlobalLock.Unlock()
	late := time.Now().Unix()
	if late-start < 3 {
		t.Error()
	}
}

func TestGTryLock(t *testing.T) {
	start := time.Now().Unix()
	c := make(chan int64)
	go func() {
		lock.GlobalLock.Lock()
		c <- time.Now().Unix()
		time.Sleep(3 * time.Second)
		defer func() {
			lock.GlobalLock.Unlock()
		}()
	}()
	<-c

	var counter int
	timer := time.After(5 * time.Second)
	for {
		ok := lock.GlobalLock.TryLock(1 * time.Second)
		counter++
		now := time.Now().Unix()

		if ok {
			if now-start < 3 {
				t.Error()
			}
			if counter < 3 {
				t.Error()
			}
			return
		} else {
			select {
			case <-timer:
				t.Errorf("Try Lock TimeOut, start: %d, now: %d, trycount: %d", start, now, counter)
				return
			default:

			}
		}
	}
}
