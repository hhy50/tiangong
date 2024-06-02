package session

import (
	"context"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

// Bridge session to client
// and client reponse to session
type Bridge interface {
	Transport(*protocol.DataPacket) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	Requests map[uint16]int
	src      net.Conn
	dst      *client.Client
}

// DirctClientBridging dispatcher to client
type DirctClientBridging struct {
	src *Session
	ctx context.Context
}

func (w *WirelessBridging) Transport(packet *protocol.DataPacket) error {
	buffer := buf.NewBuffer(packet.Len())
	defer buffer.Release()

	if err := protocol.EncodePacket(buffer, packet); err != nil {
		return err
	}

	var dial func(buf.Buffer) error
	if packet.Cmd() == protocol.Data {
		dial = w.dst.ReadFrom
	} else {
		dial = w.src.ReadFrom
	}

	if err := dial(buffer); err != nil {
		return err
	}
	return nil
}

func (w *DirctClientBridging) Transport(packet *protocol.DataPacket) error {

	switch packet.Status() {
	case protocol.New:
		// addr, port, timeout := protocol.DecodeTarget(packet.Body)

	case protocol.Active:
	case protocol.End:
	}

	return nil
}
