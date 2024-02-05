package buf_test

import (
	"github.com/haiyanghan/tiangong/common/buf"
	"reflect"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	ringbuffer := buf.NewRingBuffer()
	defer ringbuffer.Release()

	step := 256

	b := buf.NewBuffer(step)
	defer b.Release()

	// 写入4k
	for i := 0; i < 4096; i += step {
		_ = b.Clear()
		for j := 0; j < step; j++ {
			_ = buf.WriteByte(b, byte(j%256))
		}
		n, _ := ringbuffer.Write(b, step)
		if n != step {
			t.Error("Write error")
			return
		}
	}

	// 读取 3k
	buff := make([]byte, 1024)
	_, _ = ringbuffer.Read(buff)
	for i := 0; i < 3; i++ {
		bytes := make([]byte, 1024)
		read, _ := ringbuffer.Read(bytes)
		if read != 1024 || !reflect.DeepEqual(bytes, buff) {
			t.Error("Read error")
			return
		}
	}

	// 写入 1k
	if write, _ := ringbuffer.Write(buf.Wrap(buff), len(buff)); write != 1024 {
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
