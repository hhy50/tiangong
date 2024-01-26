package protocol

import (
	"tiangong/common"
	"tiangong/common/buf"
)

var (
	ClientMessageBodyLen = 16
)

// Heartbeat with client to server
type ClientMessageBody struct {
	Type      byte
	Reserved  [7]byte
	Timestamp uint64
}

func (b *ClientMessageBody) WriteTo(buffer buf.Buffer) error {
	buf.WriteByte(buffer, b.Type)
	buf.WriteBytes(buffer, b.Reserved[:])
	buf.WriteBytes(buffer, common.Uint64ToBytes(b.Timestamp))
	return nil
}

func (b *ClientMessageBody) ReadFrom(buffer buf.Buffer) error {
	return nil
}
