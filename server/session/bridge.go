package session

import (
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

// session to Bridge, one to one
type Bridge interface {
	Transport(*protocol.DataPacket, buf.Buffer) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

func (w *WirelessBridging) Transport(packet *protocol.DataPacket, buffer buf.Buffer) error {
	if err := protocol.EncodePacket(buffer, packet); err != nil {
		return err
	}
	if err := w.dst.Write(buffer); err != nil {
		return err
	}
	_ = buffer.Clear()
	return nil
}
