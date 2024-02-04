package protocol

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"io"
	"strconv"
)

const (
	PacketHeaderLen = 16
)

type PacketHeader struct {
	Len      uint16
	Rid      uint32
	Protocol byte
	Reserved [9]byte
}

func (h *PacketHeader) WriteTo(buffer buf.Buffer) error {
	if buffer.Cap() < PacketHeaderLen {
		return errors.NewError("write bytes len too short, minnum is "+strconv.Itoa(PacketHeaderLen)+"bytes", nil)
	}
	buf.WriteBytes(buffer, common.Uint16ToBytes(h.Len))
	buf.WriteBytes(buffer, common.Uint32ToBytes(h.Rid))
	buf.WriteByte(buffer, h.Protocol)
	buf.WriteBytes(buffer, h.Reserved[:])
	return nil
}

func (h *PacketHeader) ReadFrom(buffer buf.Buffer) error {
	if buffer.Len() < PacketHeaderLen {
		return errors.NewError("header([]byte) len too short, Minimum requirement "+strconv.Itoa(PacketHeaderLen)+"bytes", io.EOF)
	}
	h.Len, _ = buf.ReadUint16(buffer)
	h.Rid, _ = buf.ReadUint32(buffer)
	h.Protocol, _ = buf.ReadByte(buffer)
	{
		for range h.Reserved {
			_, _ = buf.ReadByte(buffer)
		}
	}
	return nil
}
