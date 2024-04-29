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
	"github.com/haiyanghan/tiangong/common/errors"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var (
	NoAlloc = net.IpAddress{0, 0, 0, 0}

	Default = Client{
		Name:     "Default",
		Internal: NoAlloc,
		Export:   []string{},
		conn:     nil,
	}
)

type Client struct {
	Name     string
	Internal net.IpAddress
	Export   []string

	ctx        context.Context
	auth       *protocol.ClientAuthBody
	conn       net.Conn
	lastAcTime time.Time
}

func (c *Client) Write(buffer buf.Buffer) error {
	return c.conn.ReadFrom(buffer)
}

func (c *Client) Read(buffer buf.Buffer) error {
	if err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	if _, err := buffer.Write(c.conn, buffer.Cap()); err != nil {
		return err
	}
	if buffer.Len() == 0 {
		return errors.NewError("Read empty packet, force offline", nil)
	}
	return nil
}

func (c *Client) Keepalive() {
	buffer := buf.NewRingBuffer()
	defer buffer.Release()
	defer CM.Offline(c)

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
				} else {
					runtime.Goexit()
				}
			} else {
				log.Info("Receive packet, %d bytes, cmd:[%d], from client[%s-%s]", packet.Header.Len, packet.Header.Cmd, c.Name, c.Internal)
				c.handlerPacket(packet)
			}
		}
	}
}

func (c *Client) handlerPacket(packet *protocol.Packet) {
	c.lastAcTime = time.Now()

	//TODO
}

func NewClient(ctx context.Context, conn net.Conn, cli *protocol.ClientAuthBody) Client {
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

		ctx:        context.WithParent(&ctx),
		auth:       cli,
		conn:       conn,
		lastAcTime: time.Now(),
	}
}
