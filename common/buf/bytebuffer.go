package buf

import (
	"io"
	"tiangong/common"
	"tiangong/common/errors"
)

type ByteBuffer struct {
	bytes []byte
	start int
	end   int
	len   int
}

func (b *ByteBuffer) Read(buff []byte) (int, error) {
	if b.start < b.end {
		l := common.Min(len(buff), b.Len())
		copy(buff[:l], b.bytes[b.start:b.start+l])
		b.start += l
		return l, nil
	}
	return 0, nil
}

func (b *ByteBuffer) Write(reader io.Reader, size int) (int, error) {
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

func (b *ByteBuffer) Release() {
	_ = b.Clear()
	b.bytes = nil
}

func (b *ByteBuffer) Clear() error {
	b.start = 0
	b.end = 0
	return nil
}
