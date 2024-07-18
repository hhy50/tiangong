package protocol

import (
	"fmt"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/net"
)

const (
	PacketHeaderLen = 12

	AuthRequest Cmd = iota
	AuthResponse
	HeartbeatRequest
	HeartbeatResponse
	Data
	DataResponse
)

var (
	EmptyBody        = make([]byte, 0)
	RidMask   uint32 = 0xFF_FF
)

type Cmd = byte

type Packet struct {
	Header PacketHeader
	Body   []byte
}

type PacketHeader struct {
	Len      uint16  // 2
	Rid      uint32  // 4
	Cmd      Cmd     // 1
	reserved [5]byte // 5
}

func (p Packet) Rid() uint16 {
	return uint16(p.Header.Rid & RidMask)
}

func (p Packet) Cmd() Cmd {
	return p.Header.Cmd
}

func (packet Packet) Len() int {
	return PacketHeaderLen + len(packet.Body)
}

func DecodePacket(buffer buf.Buffer, conn net.Conn, timeout time.Duration) (*Packet, error) {
	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, errors.NewError("SetReadDeadline err", err)
	}

	header, err := DecodePacketHeader(buffer, conn)
	if err != nil {
		return nil, err
	}
	var body []byte
	if header.Len > 0 {
		body, err = DecodePacketBody(buffer, int(header.Len), conn)
		if err != nil {
			return nil, err
		}
	}

	return &Packet{
		Header: *header,
		Body:   body,
	}, nil
}

func DecodePacketHeader(buffer buf.Buffer, conn net.Conn) (*PacketHeader, error) {
	if n, err := buffer.Write(conn, PacketHeaderLen); err != nil {
		return nil, err
	} else if n != PacketHeaderLen {
		return nil, fmt.Errorf("read bytes from connect too short, should read %d bytes, actual reading %d bytes", PacketHeaderLen, n)
	}

	header := &PacketHeader{}
	header.Len, _ = buf.ReadUint16(buffer)
	header.Rid, _ = buf.ReadUint32(buffer)
	header.Cmd, _ = buf.ReadByte(buffer)
	{
		for i := range header.reserved {
			header.reserved[i], _ = buf.ReadByte(buffer)
		}
	}
	return header, nil
}

func DecodePacketBody(buffer buf.Buffer, len int, conn net.Conn) ([]byte, error) {
	if n, err := buffer.Write(conn, len); err != nil {
		return nil, err
	} else if n != len {
		return nil, fmt.Errorf("read bytes from connect too short, should read %d bytes, actual reading %d bytes", len, n)
	}
	return buf.ReadAll(buffer)
}

func EncodePacket(buffer buf.Buffer, packet *Packet) error {
	if buffer.Cap() < packet.Len() {
		return fmt.Errorf("buffer.len too short, Minimum requirement %d bytes", packet.Len())
	}
	_ = buf.WriteBytes(buffer, common.Uint16ToBytes(packet.Header.Len))
	_ = buf.WriteBytes(buffer, common.Uint32ToBytes(packet.Header.Rid))
	_ = buf.WriteByte(buffer, packet.Header.Cmd)
	_ = buf.WriteBytes(buffer, packet.Header.reserved[:])
	_ = buf.WriteBytes(buffer, packet.Body)
	return nil
}
