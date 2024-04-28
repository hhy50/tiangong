package client

import (
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ClientManagerName = "ClientManager"

	// CM Instance
	CM *ClientManager = nil
)

type ClientManager struct {
	ctx context.Context

	// clients with Router feature
	clients map[net.IpAddress]*Client
	Lock    lock.Lock

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime time.Duration
}

func init() {
	component.Register(ClientManagerName, func(ctx context.Context) (component.Component, error) {
		CM = &ClientManager{
			ctx:         ctx,
			clients:     make(map[net.IpAddress]*Client, 128),
			Lock:        lock.NewLock(),
			MaxFreeTime: 3 * time.Minute,
		}
		return CM, nil
	})
}

func (cm *ClientManager) Start() error {
	cm.clients[NoAlloc] = &Default
	cm.startActiveCheck()
	return nil
}

func (cm *ClientManager) startActiveCheck() {
	go common.TimerFunc(func() {
		for _, cli := range cm.clients {
			if cli.Internal == NoAlloc {
				continue
			}
			now := time.Now()
			if cli.lastAcTime.Add(cm.MaxFreeTime).Before(now) {
				cm.Offline(cli)
				log.Warn("[%s] The client is not active within 3  minutes, force removal", cli.Name)
			}
		}
	}).Run(time.Minute)
}

func (cm *ClientManager) RegisterClient(c *Client) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()

	if _, f := cm.clients[c.Internal]; f {
		return errors.NewError("Unable to add existing client, duplicate internal ip: "+c.Internal.String(), nil)
	}
	cm.clients[c.Internal] = c
	log.Info("New client join. name: [%s], internal:[%s], export:[%s]", c.Name, c.Internal.String(), c.auth.Export)
	return nil
}

func (cm *ClientManager) GetClient(internal net.IpAddress) *Client {
	return cm.clients[internal]
}

func (cm *ClientManager) Offline(client *Client) {
	for _, cli := range cm.clients {
		if cli == client {
			cli.ctx.Cancel()
			_ = cli.conn.Close()
			break
		}
	}
}
