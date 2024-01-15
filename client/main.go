package main

import (
	"flag"
	"fmt"
	"os"
	"tun-demo/client"
	"tun-demo/model"
	"tun-demo/server"
)

func main() {
	//模式
	var mode string
	//配置信息
	config := model.Config{}

	flag.StringVar(&config.CIDR, "cidr", "172.16.0.1/24", "tun CIDR")
	flag.StringVar(&config.LocalAddr, "local", "0.0.0.0:7777", "local address")
	flag.StringVar(&config.ServerAddr, "server", "0.0.0.0:6666", "server address")
	flag.StringVar(&mode, "mode", "", "server mode")
	flag.Parse()

	if mode == "" {
		fmt.Println("启动模式(c/s)不可为空")
		os.Exit(-1)
	} else if mode == "c" {
		client.StartLinuxClientOfUdp(&config)
	} else if mode == "s" {
		server.StartLinuxServerOfUdp(&config)
	}

}
