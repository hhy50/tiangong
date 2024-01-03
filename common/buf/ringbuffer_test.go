package buf_test

import (
	"bytes"
	"testing"
	"tiangong/common/buf"
)

func TestRingBuffer(t *testing.T) {
	ringbuffer := buf.NewRingBuffer()
	step := 256

	// 写入4k
	for i := 0; i < 4096; i += step {
		b := bytes.Buffer{}
		for j := 0; j < step; j++ {
			b.WriteByte(byte(j % 256))
		}
		n, _ := ringbuffer.Write(&b)
		if n != step {
			t.Error("Write error")
			return
		}
	}

	// 读取 3k
	buff := make([]byte, 1024)
	for i := 0; i < 3; i++ {
		read, _ := ringbuffer.Read(buff)
		if read != 1024 {
			t.Error("Read error")
			return
		}
	}

	// 写入 1k
	if write, _ := ringbuffer.Write(bytes.NewBuffer(buff)); write != 1024 {
		t.Error("Write error")
		return
	}

	// 读取全部
	n := ringbuffer.Len()
	buff = make([]byte, n)
	if read, _ := ringbuffer.Read(buff); read != n || ringbuffer.Len() > 0 {
		t.Error("Read error")
		return
	}
}
