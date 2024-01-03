package buf

import (
	"io"
	"tiangong/common"
)

// RingBuffer is an implementation of Buffer using a ring buffer.
type RingBuffer struct {
	len    int
	buffer *block

	offset_w int
	offset_r int
}

func (b *RingBuffer) Release() error {
	//TODO implement me
	panic("implement me")
}

// Read reads data from the channel into the specified buffer.
func (b *RingBuffer) Read(buf []byte) (int, error) {
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
func (b *RingBuffer) Write(reader io.Reader) (int, error) {
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

	tmp := make([]byte, l)
	n, err := reader.Read(tmp)
	if err != nil {
		return n, err
	}

	off := 0
	for off < n {
		limit := common.Min(n, l)
		if end := offset + l; end > b.len {
			limit -= (end - b.len)
		}
		copy(b.buffer[offset:offset+limit], tmp[off:off+limit])
		off += limit
		l -= limit
	}
	b.offset_w += n
	return n, nil
}

func (b *RingBuffer) Len() int {
	return b.offset_w - b.offset_r
}
func Release() error {
	return nil
}
