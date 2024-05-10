package client

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

type Dispatcher interface {
	Dispatch() error
}

type DirctClientDispatcher struct {
	ctx context.Context
}

func (d *DirctClientDispatcher) Dispatch(target net.IpAddress, packet protocol.Packet) error {
	return nil
}
