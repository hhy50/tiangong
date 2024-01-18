package main

import (
	"flag"
	"tiangong/common"
	"tiangong/common/log"
)

var (
	Server  string
	Port    int
	Token   string
	SubHost string
)

func init() {
	// TODO parse config file
	flag.StringVar(&Server, "server", "127.0.0.1", "Specify target server")
	flag.IntVar(&Port, "port", 2024, "Specify target port")
	flag.StringVar(&Token, "token", "tiangong", "Token")
	flag.StringVar(&SubHost, "subHost", "", "SubHost")
}

func main() {
	log.InitLog()
	processor := NewProcessor()
	if err := processor.Start(); err != nil {
		log.Error("start fail, ", err)
		return
	}
	log.Info("Kernel client start success")
	<-common.WaitSignal()
}
