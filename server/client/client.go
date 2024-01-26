package client

import (
	"context"
	"runtime"
	"tiangong/common/buf"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

type Client struct {
	Name     string
	Internal net.IpAddress

	ctx        context.Context
	cancel     context.CancelFunc
	auth       *protocol.ClientAuth
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
	return nil
}

func (c *Client) Keepalive() {
	buffer := buf.NewRingBuffer()
	defer buffer.Release()
	defer c.Offline()

	select {
	case <-c.ctx.Done():
		runtime.Goexit()
	default:
		if err := c.Read(buffer); err != nil {
			log.Error("read bytes from client error, ", err)
			return
		}
		handlerResponse(buffer)
	}

}

func handlerResponse(buffer buf.Buffer) {
	//TODO
	buffer.Clear()
}

func (c *Client) Offline() {
	_ = c.conn.Close()
	c.cancel()
	delete(Clients, c.Internal)
	delete(ClientNames, c.Name)
	log.Warn("cliet [%s-%s] offlined...", c.Name, c.Internal.String())
}
