package protocol

import (
	"io"
	"strconv"
	"tiangong/common"
	"tiangong/common/errors"
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

func (h *PacketHeader) Unmarshal(header []byte) error {
	if len(header) < PacketHeaderLen {
		return errors.NewError("header([]byte) len too short, Minimum requirement "+strconv.Itoa(PacketHeaderLen)+"bytes", io.EOF)
	}
	h.Len = common.Uint16(header[0:2])
	h.Rid = common.Uint32(header[2:6])
	h.Protocol = header[6]
	return nil
}
