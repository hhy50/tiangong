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

	// Clients with Router feature
	Clients = make(map[net.IpAddress]*Client, 128)
	Lock    = lock.NewLock()

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime = 3 * time.Minute

	safe = Client{
		Name:     "Instance",
		Internal: NoAlloc,
		Export:   []string{},
	}
)

type ClientManager struct {
	ctx context.Context

	// Clients with Router feature
	clients map[net.IpAddress]*Client
	Lock    lock.Lock

	// MaxFreeTime The maximum idle time allowed to the client
	MaxFreeTime time.Duration
}

func init() {
	component.Register(ClientManagerName, func(ctx context.Context) (component.Component, error) {
		return &ClientManager{
			ctx:         ctx,
			clients:     make(map[net.IpAddress]*Client, 128),
			Lock:        lock.NewLock(),
			MaxFreeTime: 3 * time.Minute,
		}, nil
	})
}

func GetClient(internal net.IpAddress) *Client {
	return Clients[internal]
}

func (manager *ClientManager) Start() error {
	RegistClient(&safe)

	manager.startActiveCheck()
	return nil
}

func (manager ClientManager) startActiveCheck() {
	go common.TimerFunc(func() {
		for _, cli := range Clients {
			if cli.Internal == NoAlloc {
				continue
			}
			now := time.Now()
			if cli.lastAcTime.Add(MaxFreeTime).Before(now) {
				cli.Offline()
				log.Warn("[%s] The client is not active within 3  minutes, force removal", cli.GetName())
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
	Clients[c.Internal] = c
	log.Info("New client join. name: [%s], internal:[%s]", c.Name, c.Internal.String())
	return nil
}
