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

// func NewDataPacketW(host string, port, timeout, rid uint16) *DataPacket {
// 	buffer := buf.NewBuffer(256)
// 	defer buffer.Release()

// 	buf.WriteBytes(buffer, []byte(host))
// 	buf.WriteBytes(buffer, common.Uint16ToBytes(port))
// 	buf.WriteBytes(buffer, common.Uint16ToBytes(timeout))

// 	body, _ := buf.ReadAll(buffer)

// 	packet := DataPacket{
// 		Header: PacketHeader{
// 			Len:      uint16(buffer.Len()),
// 			Rid:      rid,
// 			Cmd:      Data,
// 			reserved: [5]byte{New},
// 		},
// 		Body: body,
// 	}
// 	return &packet
// }

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
