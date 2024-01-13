package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"tiangong/common/log"
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

func main() {
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
