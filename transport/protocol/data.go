package protocol

import "github.com/haiyanghan/tiangong/common"

type Status = byte

const (
	New Status = iota
	Active
	End
)

type DataPacket = Packet

func (packet *DataPacket) Status() Status {
	return packet.Header.reserved[0]
}

func NewDataPacket(rid uint16, status Status, body []byte) *DataPacket {
	packet := DataPacket{
		Header: PacketHeader{
			Len:      uint16(len(body)),
			Rid:      rid,
			Cmd:      Data,
			reserved: [5]byte{status},
		},
		Body: body,
	}
	return &packet
}

func EncodeTarget(addr string, port uint16, timeout uint16) []byte {
	n := len(addr)
	bytes := make([]byte, n+4)
	copy(bytes[:n], []byte(addr))
	copy(bytes[n:n+2], common.Uint16ToBytes(port))
	copy(bytes[n+2:n+4], common.Uint16ToBytes(timeout))
	return bytes
}

func DecodeTarget(bytes []byte) (addr string, port uint16, timeout uint16) {
	n := len(bytes)
	timeout = common.Uint16(bytes[n-2:n])
	port = common.Uint16(bytes[n-4:n-2])
	addr = common.String(bytes[0:n-4])
	return
}
