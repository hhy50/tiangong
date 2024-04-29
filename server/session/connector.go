package session

import (
	"fmt"
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/auth"
	"github.com/haiyanghan/tiangong/server/client"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	ConnectorCompName = "SessionConnect"
)

type Connector struct {
	net.TcpServer
}

type Config struct {
	Port int `prop:"port" default:"2024"`
}

func init() {
	component.Register(ConnectorCompName, func(ctx context.Context) (component.Component, error) {
		config := Config{}
		err := conf.LoadConfig("session", &config)
		if err != nil {
			return nil, err
		}
		host := conf.GetOrDefault("server.host", "127.0.0.1").(string)
		return &Connector{
			TcpServer: net.NewTcpServer(host, config.Port, ctx),
		}, nil
	})
}

func (c *Connector) Start() error {
	return c.TcpServer.ListenTCP(connHandler)
}

func connHandler(ctx context.Context, conn net.Conn) error {
	_, sessionAuth, err := auth.AuthToken(conn)
	if err != nil {
		return err
	}

	subHost := net.ParseFromStr(sessionAuth.SubHost)
	cm := ctx.Value(client.ClientManagerName).(*client.ClientManager)
	dstClient := cm.GetClient(subHost)

	if dstClient == nil {
		return errors.NewError(fmt.Sprintf("subhost '%s' not fount", subHost.String()), nil)
	}

	s := NewSession(sessionAuth.Token, conn, ctx, dstClient)
	manager := ctx.Value(ManagerCompName).(*Manager)
	if err = manager.AddSession(subHost, &s); err != nil {
		return err
	}
	go s.Work()
	return nil
}
