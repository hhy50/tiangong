package buf

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/lock"
	"io"
	"sync"
)

// RingBuffer is an implementation of Buffer using a ring buffer.
type RingBuffer struct {
	len    int
	buffer *block

	offset_w int
	offset_r int
	lock     lock.Rwlock
	once     sync.Once
}

func (b *RingBuffer) Clear() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.offset_w, b.offset_r = 0, 0
	return nil
}

func (b *RingBuffer) Release() {
	b.once.Do(func() {
		b.Clear()
		b.buffer = nil
	})
}

// Read reads data from the channel into the specified buffer.
func (b *RingBuffer) Read(buf []byte) (int, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.offset_r < b.offset_w {
		l := len(buf)
		off := 0

		// Calculate relative offset
		r := b.offset_r & (b.len - 1)
		// w = b.offset_w & (b.len - 1)

		// Determine the length to read
		rl := common.Min(b.offset_w-b.offset_r, l)
		for off < rl {
			cr := common.Min(rl, b.len-r)
			copy(buf[off:], b.buffer[r:r+cr])
			off += cr
			r = 0
		}
		b.offset_r += rl
		return rl, nil
	}
	return 0, nil
}

// Write writes data from the specified buffer into the channel.
func (b *RingBuffer) Write(reader io.Reader, size int) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	r := b.offset_r & (b.len - 1)
	w := b.offset_w & (b.len - 1)

	offset, l := func() (int, int) {
		if b.offset_r == b.offset_w {
			return 0, b.len
		}
		if (b.offset_w-b.offset_r)&(b.len-1) == 0 {
			return b.len, 0
		}
		if r > w {
			return w, r - w
		}
		return w, b.len - w + r
	}()

	if l < size {
		return 0, NoSpace
	}

	tmp := make([]byte, size)
	n, err := reader.Read(tmp)
	if err != nil {
		return n, err
	}

	off := 0
	limit := n
	for off < n {
		if end := offset + limit; end > b.len {
			limit = b.len - offset
		}
		copy(b.buffer[offset:offset+limit], tmp[off:off+limit])
		off += limit
	}
	b.offset_w += n
	return n, nil
}

func (b *RingBuffer) Len() int {
	return b.offset_w - b.offset_r
}

func (b *RingBuffer) Cap() int {
	if b.offset_r == b.offset_w {
		return b.len
	}
	if (b.offset_w-b.offset_r)&(b.len-1) == 0 {
		return 0
	}
	r := b.offset_r & (b.len - 1)
	w := b.offset_w & (b.len - 1)
	if r > w {
		return r - w
	}
	return b.len - w + r
}
