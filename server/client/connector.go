package client

import (
	"github.com/haiyanghan/tiangong/common/context"

	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/auth"
	"github.com/haiyanghan/tiangong/server/component"
)

type Config struct {
	Key  string `prop:"key" default:""`
	Port int    `prop:"port" default:"2024"`
}

type ClientConnentor struct {
	ctx    context.Context
	tcpSrv net.TcpServer
}

var (
	Conf Config
)

func init() {
	component.Register("ClientConnentor", func(ctx context.Context) (component.Component, error) {
		err := conf.LoadConfig("client", &Conf)
		if err != nil {
			return nil, err
		}
		host := conf.GetOrDefault("server.host", "127.0.0.1").(string)
		return &ClientConnentor{
			ctx:    ctx,
			tcpSrv: net.NewTcpServer(host, Conf.Port, ctx),
		}, nil
	})
}

func (client *ClientConnentor) Start() error {
	return client.tcpSrv.ListenTCP(connHandler)
}

func connHandler(ctx context.Context, conn net.Conn) error {
	_, cli, err := auth.AuthKey(Conf.Key, conn)
	if err != nil {
		return err
	}

	c := NewClient(ctx, conn, cli)
	if err := RegistClient(&c); err != nil {
		return err
	}
	go c.Keepalive()
	return nil
}
