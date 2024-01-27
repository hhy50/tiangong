package protocol

import (
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"io"
	"strconv"
	"unsafe"
)

const (
	PacketHeaderLen = int(unsafe.Sizeof((*PacketHeader)(nil)))
)

type PacketHeader struct {
	Len      uint16
	Rid      uint32
	Protocol byte
}

func (h *PacketHeader) WriteTo(buffer buf.Buffer) error {
	if buffer.Cap() < PacketHeaderLen {
		return errors.NewError("write bytes len too short, minnum is "+strconv.Itoa(PacketHeaderLen)+"bytes", nil)
	}

	return nil
}

func (h *PacketHeader) ReadFrom(buffer buf.Buffer) error {
	if buffer.Len() < PacketHeaderLen {
		return errors.NewError("header([]byte) len too short, Minimum requirement "+strconv.Itoa(PacketHeaderLen)+"bytes", io.EOF)
	}
	h.Len, _ = buf.ReadUint16(buffer)
	h.Rid, _ = buf.ReadUint32(buffer)
	h.Protocol, _ = buf.ReadByte(buffer)
	return nil
}
