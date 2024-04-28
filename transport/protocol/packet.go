package protocol

import (
	"io"
	"strconv"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
)

const (
	PacketHeaderLen = 10

	AuthRequest Cmd = iota
	AuthResponse
	HeartbeatRequest
	HeartbeatResponse
	Data
)

type Cmd = byte

type PacketHeader struct { // 10
	Len      uint16  // 2
	Rid      uint16  // 2
	Cmd      Cmd     // 1
	reserved [5]byte // 5
}

func (h *PacketHeader) WriteTo(buffer buf.Buffer) error {
	if buffer.Cap() < PacketHeaderLen {
		return errors.NewError("write bytes len too short, minimum is "+strconv.Itoa(PacketHeaderLen)+"bytes", nil)
	}
	_ = buf.WriteBytes(buffer, common.Uint16ToBytes(h.Len))
	_ = buf.WriteBytes(buffer, common.Uint16ToBytes(h.Rid))
	_ = buf.WriteByte(buffer, h.Cmd)
	_ = buf.WriteBytes(buffer, h.reserved[:])
	return nil
}

func (h *PacketHeader) ReadFrom(buffer buf.Buffer) error {
	if buffer.Len() < PacketHeaderLen {
		return errors.NewError("header([]byte) len too short, Minimum requirement "+strconv.Itoa(PacketHeaderLen)+"bytes", io.EOF)
	}
	h.Len, _ = buf.ReadUint16(buffer)
	h.Rid, _ = buf.ReadUint16(buffer)
	h.Cmd, _ = buf.ReadByte(buffer)
	{
		for i := range h.reserved {
			h.reserved[i], _ = buf.ReadByte(buffer)
		}
	}
	return nil
}
