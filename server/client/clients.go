package client

import (
	"context"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common"

	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/transport/protocol"
)

var (
	// Clients with Router feature
	Clients     = make(map[net.IpAddress]*Client, 128)
	ClientNames = make(map[string]*Client, 128)
	Lock        = lock.NewLock()

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime = 3 * time.Minute
)

func StartActiveCheck(ctx context.Context) {
	go common.TimerFunc(func() {
		now := time.Now()
		select {
		case <-ctx.Done():
			runtime.Goexit()
		default:
			for _, cli := range Clients {
				if cli.lastAcTime.Add(MaxFreeTime).Before(now) {
					cli.Offline()
					log.Warn("[%s] The client is not active within 3  minutes, force removal", cli.GetName())
				}
			}
		}
	}).Run(time.Minute)
}

func RegistClient(c *Client) error {
	Lock.Lock()
	defer Lock.Unlock()

	if _, f := Clients[c.Internal]; f {
		return errors.NewError("Unable to add existing client, duplicate internal ip: "+c.Internal.String(), nil)
	}
	if _, f := ClientNames[c.Name]; f {
		return errors.NewError("Unable to add existing client, duplicate name: "+c.Name, nil)
	}
	Clients[c.Internal] = c
	ClientNames[c.Name] = c
	log.Info("New client join. name: [%s], internal:[%s]", c.Name, c.Internal.String())
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
