package buf

import "io"

const (
	BlockLen = 4096
)

// 4k bytes memory block
type block [BlockLen]byte

type Buffer interface {
	Read([]byte) (int, error)
	Write(reader io.Reader, len int) (int, error)
	Len() int
	Release()
	Clear() error
}

func NewRingBuffer() Buffer {
	return &RingBuffer{
		len:    BlockLen,
		buffer: &block{},
	}
}

func Wrap(bytes []byte) Buffer {
	// TODO
	return nil
}
