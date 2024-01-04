package buf

import "io"

const (
	BlockLen = 4096
)

// 4k bytes memory block
type block [BlockLen]byte

type Buffer interface {
	Read([]byte) (int, error)
	Write(reader io.Reader) (int, error)
	Len() int
	Release() error
}

func NewRingBuffer() Buffer {
	return &RingBuffer{
		len:    BlockLen,
		buffer: &block{},
	}
}

func NewBuffer(len int) Buffer {
	return nil
}
