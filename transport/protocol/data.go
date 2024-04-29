package protocol

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
