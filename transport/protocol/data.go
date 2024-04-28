package protocol

type Status = byte

const (
	New Status = iota
	Active
	End
)

type DataPacket = PacketHeader

func (packet *DataPacket) Status() Status {
	return packet.reserved[0]
}
