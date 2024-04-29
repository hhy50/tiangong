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
	Name     string
	Internal net.IpAddress
	Export   []string

	conn       net.Conn
	ctx        context.Context
	lastAcTime time.Time
}

func (c *Client) Write(buffer buf.Buffer) error {
	return c.conn.ReadFrom(buffer)
}

func (c *Client) Keepalive() {
	buffer := buf.NewRingBuffer()
	defer func() {
		_ = c.conn.Close()
		buffer.Release()

		cm := c.ctx.Value(ManagerName).(*Manager)
		cm.Offline(c)
	}()

	for {
		select {
		case <-c.ctx.Done():
			runtime.Goexit()
		default:
			if err := c.conn.SetReadDeadline(time.Now().Add(time.Minute)); err != nil {
				return
			}
			if packet, err := protocol.DecodePacket(buffer, c.conn); err != nil {
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

	// TODO
}

func NewClient(ctx context.Context, cli *protocol.ClientAuthBody) Client {
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
		Name:     cli.Name,
		Internal: getInternalIpFromReq(),
		Export:   strings.Split(cli.Export, ","),

		conn:       ctx.Value(net.ConnValName).(net.Conn),
		ctx:        ctx,
		lastAcTime: time.Now(),
	}
}
