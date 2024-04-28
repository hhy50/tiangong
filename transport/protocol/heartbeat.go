package protocol

const (
	HeartbeatPacketLen = PacketHeaderLen
)

type HeartbeatPacketHeader = PacketHeader

func NewHeartbeatPacket() *HeartbeatPacketHeader {
	return &HeartbeatPacketHeader{
		Rid:      0,
		Len:      uint16(HeartbeatPacketLen),
		Cmd:      HeartbeatRequest,
		reserved: [5]byte{},
	}
}
