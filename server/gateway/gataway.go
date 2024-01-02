package gateway

import "tiangong/kernel/transport/protocol"

type Gateway interface {
	Inbound() error
	Outbound() error
}

type Handler func(header *protocol.RequestHeader) *protocol.ResponseBody

func (f *Handler) Inbound() {

}

func (f *Handler) Outbound() {

}
