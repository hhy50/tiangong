package client

import (
	"fmt"
	"runtime"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/lock"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ManagerName = "Manager"
)

type Manager struct {
	ctx context.Context

	// clients with Router feature
	clients map[net.IpAddress]*Client
	Lock    lock.Lock

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime time.Duration
}

func init() {
	component.Register(ManagerName, func(ctx context.Context) (component.Component, error) {
		return &Manager{
			ctx:         ctx,
			clients:     make(map[net.IpAddress]*Client, 128),
			Lock:        lock.NewLock(),
			MaxFreeTime: 3 * time.Minute,
		}, nil
	})
}

func (cm *Manager) Start() error {
	defaultClient := Client{
		Name:     "Default",
		Internal: NoAlloc,
		Export:   []string{},
		ctx:      cm.ctx,
	}

	cm.RegisterClient(&defaultClient)
	cm.startActiveCheck()
	return nil
}

func (cm *Manager) startActiveCheck() {
	go common.TimerFunc(func() {
		for _, cli := range cm.clients {
			if cli.Internal == NoAlloc {
				continue
			}
			now := time.Now()
			if cli.lastAcTime.Add(cm.MaxFreeTime).Before(now) {
				cm.Offline(cli)
				log.Warn("[%s] The client is not active within 3 minutes, force removal", cli.Name)
			}
		}
		runtime.GC()
	}).Run(time.Minute)
}

func (cm *Manager) GetClient(internal net.IpAddress) *Client {
	return cm.clients[internal]
}

func (cm *Manager) RegisterClient(c *Client) error {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()

	if _, f := cm.clients[c.Internal]; f {
		return fmt.Errorf("unable to add existing client, duplicate internal ip: %s", c.Internal.String())
	}
	cm.clients[c.Internal] = c
	log.Info("New client join. name: [%s], internal:[%s], export:[%+v]", c.Name, c.Internal.String(), c.Export)
	return nil
}

func (cm *Manager) Offline(client *Client) {
	if cli, f := cm.clients[client.Internal]; f {
		if cli == client {
			cli.ctx.Cancel()
			log.Warn("Client [%s-%s] offline...", cli.Name, cli.Internal.String())
		}
		delete(cm.clients, client.Internal)
	}
}
