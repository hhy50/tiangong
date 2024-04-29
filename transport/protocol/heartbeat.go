package protocol

const (
	HeartbeatPacketLen = PacketHeaderLen
)

func NewHeartbeatPacket() *Packet {
	return &Packet{
		Header: PacketHeader{
			Rid:      0,
			Len:      uint16(HeartbeatPacketLen),
			Cmd:      HeartbeatRequest,
			reserved: [5]byte{},
		},
	}
}
