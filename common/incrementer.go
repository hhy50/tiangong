package common

import (
	"sync/atomic"
)

type Range = [2]uint64

type Incrementer struct {
	Range Range

	counter uint64
}

func (i *Incrementer) Next() uint64 {
	min := i.Range[0]
	max := i.Range[1]
	n := atomic.AddUint64(&i.counter, 1)
	return (n % max) + min
}
