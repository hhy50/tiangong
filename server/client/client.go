package client

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common/errors"

	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
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
	if buffer.Len() == 0 {
		c.Offline()
		return errors.NewError("read empty packet, force offline", nil)
	}
	return nil
}

func (c *Client) Keepalive() {
	buffer := buf.NewRingBuffer()
	defer buffer.Release()
	defer c.Offline()
	for {
		select {
		case <-c.ctx.Done():
			runtime.Goexit()
		default:
			if err := c.Read(buffer); err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				} else {
					runtime.Goexit()
				}
			}
			c.lastAcTime = time.Now()
			log.Debug("Receive %d bytes from client[%s]", buffer.Len(), c.GetName())
			handlerResponse(buffer)
		}
	}
}

func handlerResponse(buffer buf.Buffer) {
	//TODO
	buffer.Clear()
}

func (c *Client) Offline() {
	Lock.Lock()
	defer Lock.Unlock()

	_ = c.conn.Close()
	c.cancel()
	delete(Clients, c.Internal)
	delete(ClientNames, c.Name)
	log.Warn("Client [%s] is offlined...", c.GetName())
}

func (c *Client) GetName() string {
	return fmt.Sprintf("%s-%s", c.Name, c.Internal.String())
}
