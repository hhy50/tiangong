package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tiangong/common/log"
	"tiangong/kernel/client"
	"tiangong/kernel/model"
)

var (
	server  string
	port    int
	token   string
	subHost string
)

func init() {
	// TODO parse config file
	flag.StringVar(&server, "server", "127.0.0.1", "Specify target server")
	flag.IntVar(&port, "port", 2024, "Specify target port")
	flag.StringVar(&token, "token", "tiangong", "Token")
	flag.StringVar(&subHost, "subHost", "", "SubHost")
}

func main1() {
	log.InitLog()
	processor := NewProcessor(server, port, token, subHost)
	if err := processor.Start(); err != nil {
		log.Error("start fail, ", err)
		return
	}
	log.Info("Kernel client start success")
	pauseProcess()
}

func pauseProcess() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
}

func main() {
	var mode string
	//配置信息
	config := model.Config{}

	flag.StringVar(&config.CIDR, "cidr", "172.16.0.1/24", "tun CIDR")
	flag.StringVar(&config.LocalAddr, "localAddr", "0.0.0.0:7777", "local address")
	flag.StringVar(&config.ServerAddr, "serverAddr", "0.0.0.0:6666", "server address")
	flag.StringVar(&mode, "mode", "", "server mode")
	flag.Parse()

	if mode == "" {
		fmt.Println("启动模式(c/s)不可为空")
		os.Exit(-1)
	} else if mode == "c" {
		client.StartWinClient(&config)
	} else if mode == "s" {
		//server.StartLinuxServerOfUdp(&config)
	}
}
