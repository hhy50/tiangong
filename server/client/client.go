package client

import (
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/server/internal"

	"github.com/google/uuid"
	"github.com/haiyanghan/tiangong/common"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var (
	NoAlloc = net.IpAddress{0, 0, 0, 0}
)

type Client struct {
	net.Conn
	Name     string
	Internal net.IpAddress
	Export   []string

	ctx        context.Context
	lastAcTime time.Time
}

func (c *Client) Keepalive() {
	buffer := buf.NewRingBuffer()
	defer func() {
		_ = c.Conn.Close()
		buffer.Release()

		cm := c.ctx.Value(ManagerName).(*Manager)
		cm.Offline(c)
	}()

	for {
		select {
		case <-c.ctx.Done():
			runtime.Goexit()
		default:
			if packet, err := protocol.DecodePacket(buffer, c.Conn, 15*time.Second); err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return
			} else {
				log.Info("Receive packet, %d bytes, cmd:[%d], from client[%s-%s]", packet.Header.Len, packet.Header.Cmd, c.Name, c.Internal)
				c.handlerPacket(packet)
			}
		}
	}
}

func (c *Client) handlerPacket(packet *protocol.Packet) {
	c.lastAcTime = time.Now()

	switch packet.Cmd() {
	case protocol.DataResponse:

	case protocol.HeartbeatRequest:

	}
}

func NewClient(ctx context.Context, cli *protocol.ClientAuthBody) Client {
	ctx = context.WithParent(ctx)
	getInternalIpFromReq := func() net.IpAddress {
		if len(cli.Internal) == 4 || reflect.DeepEqual(cli.Internal, NoAlloc) {
			i := cli.Internal
			return net.IpAddress{i[0], i[1], i[2], i[3]}
		}
		return internal.GeneraInternalIp()
	}

	if common.IsEmpty(cli.Name) {
		uid, _ := uuid.NewUUID()
		cli.Name = uid.String()
	}

	return Client{
		Conn:     ctx.Value(net.ConnValName).(net.Conn),
		Name:     cli.Name,
		Internal: getInternalIpFromReq(),
		Export:   strings.Split(cli.Export, ","),

		ctx:        ctx,
		lastAcTime: time.Now(),
	}
}
