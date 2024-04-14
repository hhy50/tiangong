package protocol

import (
	"io"
	"strconv"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/net"
)

const (
	PacketHeaderLen = 20
)

type PacketHeader struct { // 20
	Len      uint16        // 2
	Rid      uint32        // 4
	Protocol byte          // 1
	Host     net.IpAddress // 4
	Port     net.Port      // 2
	Reserved [6]byte       // 6
	Status   byte          // 1
}

// func (h *PacketHeader) WriteTo(buffer buf.Buffer) error {
// 	if buffer.Cap() < PacketHeaderLen {
// 		return errors.NewError("write bytes len too short, minnum is "+strconv.Itoa(PacketHeaderLen)+"bytes", nil)
// 	}
// 	buf.WriteBytes(buffer, common.Uint16ToBytes(h.Len))
// 	buf.WriteBytes(buffer, common.Uint32ToBytes(h.Rid))
// 	buf.WriteByte(buffer, h.Protocol)
// 	buf.WriteByte(buffer, h.HostLen)
// 	buf.WriteBytes(buffer, common.Uint16ToBytes(h.Port))
// 	buf.WriteBytes(buffer, h.Reserved[:])
// 	return nil
// }

func (h *PacketHeader) ReadFrom(buffer buf.Buffer) error {
	if buffer.Len() < PacketHeaderLen {
		return errors.NewError("header([]byte) len too short, Minimum requirement "+strconv.Itoa(PacketHeaderLen)+"bytes", io.EOF)
	}
	h.Len, _ = buf.ReadUint16(buffer)
	h.Rid, _ = buf.ReadUint32(buffer)
	h.Protocol, _ = buf.ReadByte(buffer)
	{
		hostPort, _ := buf.ReadBytes(buffer, 6)
		h.Host = net.ParseFromBytes(hostPort[:5])
		h.Port = net.Port(common.Uint16(hostPort[5:]))
	}
	{
		for range h.Reserved {
			_, _ = buf.ReadByte(buffer)
		}
	}
	h.Status, _ = buf.ReadByte(buffer)
	return nil
}
