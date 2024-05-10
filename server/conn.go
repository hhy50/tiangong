package server

import (
	"fmt"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

func ConnHandler(ctx context.Context, conn net.Conn) error {
	ctx = context.WithParent(ctx)
	ctx.AddValue(net.ConnValName, conn)

	buffer := buf.NewBuffer(4096)
	defer buffer.Release()

	if packet, err := protocol.DecodePacket(buffer, conn, Timout); err != nil {
		return err
	} else {
		switch packet.AuthType() {
		case protocol.AuthClient:
			return AuthKey(packet, ctx)
		case protocol.AuthSession:
			return AuthToken(packet, ctx)
		default:
			return fmt.Errorf("unsupported connect type, type=%d", packet.AuthType())
		}
	}
}
