package buf

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/lock"
	"io"
	"sync"
)

type ByteBuffer struct {
	bytes []byte
	start int
	end   int
	len   int

	lock lock.Rwlock
	once sync.Once
}

func (b *ByteBuffer) Read(buff []byte) (int, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.start < b.end {
		l := common.Min(len(buff), b.Len())
		copy(buff[:l], b.bytes[b.start:b.start+l])
		b.start += l
		return l, nil
	}
	return 0, nil
}

func (b *ByteBuffer) Write(reader io.Reader, size int) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.end < b.len {
		l := b.len - b.end
		if l < size {
			return 0, NoSpace
		}
		n, err := reader.Read(b.bytes[b.end : b.end+size])
		b.end += n
		return n, err
	}
	return 0, errors.NewError("current buffer Unable to write, no space", nil)
}

func (b *ByteBuffer) Len() int {
	return b.end - b.start
}

func (b *ByteBuffer) Cap() int {
	return b.len - b.end
}

func (b *ByteBuffer) Release() {
	_ = b.Clear()
	b.bytes = nil
}

func (b *ByteBuffer) Clear() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.start = 0
	b.end = 0
	return nil
}
