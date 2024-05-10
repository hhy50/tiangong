package session

import (
	"context"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

// session to Bridge, one to one
type Bridge interface {
	Transport(*protocol.DataPacket) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

// DirctClientBridging dispatcher to client
type DirctClientBridging struct {
	ctx      context.Context
	Requests map[uint16]int
}

func (w *WirelessBridging) Transport(packet *protocol.DataPacket) error {
	buffer := buf.NewBuffer(packet.Len())
	defer buffer.Release()

	if err := protocol.EncodePacket(buffer, packet); err != nil {
		return err
	}
	if err := w.dst.Write(buffer); err != nil {
		return err
	}
	return nil
}

func (w *DirctClientBridging) Transport(packet *protocol.DataPacket) error {
	return nil
}
