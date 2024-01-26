package client

import (
	"context"
	"tiangong/common/errors"
	"tiangong/common/lock"
	"tiangong/common/log"
	"tiangong/common/net"
	"tiangong/kernel/transport/protocol"
	"time"
)

var (
	// Clients with Router feature
	Clients     = make(map[net.IpAddress]*Client, 128)
	ClientNames = make(map[string]*Client, 128)
	Lock        = lock.NewLock()

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime = 10 * time.Minute
)

func init() {
	// heartbeat check
	go func() {
		ticker := time.NewTicker(time.Minute)
		for {
			<-ticker.C
			log.Debug("heartbeat check...")
			now := time.Now()
			for addr, cli := range Clients {
				if cli.lastAcTime.Add(MaxFreeTime).Before(now) {
					cli.Offline()
					log.Warn("[%s-%s] The client is not active within 10 minutes, force removal", addr.String(), cli.Name)
				}
			}
		}
	}()
}

func RegistClient(c *Client) error {
	Lock.Lock()
	defer Lock.Unlock()

	if _, f := Clients[c.Internal]; f {
		return errors.NewError("Unable to add existing client, name: "+c.Internal.String(), nil)
	}
	if _, f := ClientNames[c.Name]; f {
		return errors.NewError("Unable to add existing client, name: "+c.Name, nil)
	}
	Clients[c.Internal] = c
	ClientNames[c.Name] = c
	return nil
}

func NewClient(ctx context.Context, internalIP net.IpAddress, cli *protocol.ClientAuth, conn net.Conn) Client {
	ctx, cancel := context.WithCancel(ctx)

	return Client{
		Name:       cli.Name,
		Internal:   internalIP,
		ctx:        ctx,
		cancel:     cancel,
		auth:       cli,
		conn:       conn,
		lastAcTime: time.Now(),
	}
}
