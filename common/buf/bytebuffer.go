package buf

import (
	"io"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/lock"
)

type ByteBuffer struct {
	bytes []byte
	start int
	end   int
	len   int

	lock lock.Rwlock
}

func (b *ByteBuffer) Read(buff []byte) (int, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.start < b.end {
		ln := common.Min(len(buff), b.end-b.start)
		copy(buff[:ln], b.bytes[b.start:b.start+ln])
		b.start += ln
		return ln, nil
	}
	return 0, nil
}

func (b *ByteBuffer) Write(reader io.Reader, size int) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.end < b.len {
		if b.end+size > b.len {
			return 0, NoSpace
		}
		n, err := reader.Read(b.bytes[b.end : b.end+size])
		b.end += n
		return n, err
	}
	return 0, NoSpace
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
