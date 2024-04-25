package session

import (
	"net/http"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var HTTP_CLIENT = &http.Client{}

// session to Bridge, one to one
type Bridge interface {
	Transport(*protocol.PacketHeader, buf.Buffer) error
}

// WirelessBridging point to point
type WirelessBridging struct {
	dst *client.Client
}

func (w *WirelessBridging) Transport(h *protocol.PacketHeader, buffer buf.Buffer) error {
	switch h.Status {
	case protocol.New:
		w.dst.WriteHeader(h)
		w.dst.WriteBody(buffer)
	case protocol.Active:
		w.dst.WriteBody(buffer)
	case protocol.End:
		w.dst.WriteHeader(h)
	}
	return nil
}
