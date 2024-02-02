package session

import (
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Bridge interface {
	Transport(protocol.PacketHeader, buf.Buffer) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

func (w *WirelessBridging) Transport(h protocol.PacketHeader, buffer buf.Buffer) error {
	_, _ = buf.ReadAll(buffer)
	// log.Info("[%s]", bytes)
	return nil
}
