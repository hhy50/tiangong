package server

import (
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/server/component"
)

var (
	Config = struct {
		Host string `prop:"host"`
		Port int    `prop:"port"`
		Key  string `prop:"key"`
	}{}
)

func init() {
	component.Register("TcpServer", func(ctx context.Context) (component.Component, error) {
		conf.LoadConfig("server", &Config)
		return component.FuncComponent(func() error {
			tcpServer := net.NewTcpServer(Config.Host, Config.Port, ctx)
			return tcpServer.ListenTCP(ConnHandler)
		}), nil
	})
}
