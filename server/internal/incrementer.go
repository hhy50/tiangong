package internal

import (
	"sync/atomic"
)

type Range = [2]int

type Incrementer struct {
	count int32
	Range Range
}

func (i *Incrementer) Next() int {
	min := i.Range[0]
	max := i.Range[1]
	n := atomic.AddInt32(&i.count, 1)
	return (int(n) % max) + min
}
