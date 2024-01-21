package session

import (
	"tiangong/common/buf"
	"tiangong/common/log"
	"tiangong/kernel/transport/protocol"
	"tiangong/server/client"
)

type Bridge interface {
	Transport(protocol.PacketHeader, buf.Buffer) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

func (w *WirelessBridging) Transport(protocol.PacketHeader, buf.Buffer) error {
	log.Info("aaa")
	return nil
}
