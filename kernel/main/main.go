package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/common/net"
	"github.com/haiyanghan/tiangong/kernel/proxy"
)

var (
	cp               string
	ConnTimeout      = 30 * time.Second
	HandshakeTimeout = ConnTimeout
	Config           = config{}
)

func init() {
	flag.StringVar(&cp, "conf", "", "config file path")
}

func main() {
	flag.Parse()
	log.InitLog()

	conf.LoadConfig(cp, &Config, conf.EmptyDefaultValueFunc)

	proxyServer := net.NewTcpServer(Config.ProxyHost, Config.ProxyPort, context.Background())
	proxyServer.ListenTCP(StartListener)

	if err := proxy.SetProxy(GetProxyAddr(), []string{"192.168.110.*"}); err != nil {
		log.Warn("Set System proxy error", err)
	}
	defer proxy.ResetProxy()

	log.Info("Kernel client start success")
	<-common.WaitSignal()
}

func GetProxyAddr() string {
	return fmt.Sprintf("%s:%d", Config.ProxyHost, Config.ProxyPort)
}
