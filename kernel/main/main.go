package main

import (
	"flag"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/buf"
	"github.com/haiyanghan/tiangong/common/log"
	"os"
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
	flag.Parse()
	log.InitLog()
	processor := NewProcessor()
	if err := processor.Start(); err != nil {
		log.Error("kernel start error, ", err)
		return
	}
	log.Info("Kernel client start success")
	go func() {
		log.Info("opening stdin io...")
		for {
			buffer := buf.NewRingBuffer()
			if _, err := buffer.Write(os.Stdin, buffer.Cap()); err != nil {
				log.Error("Read form stdio error", err)
			} else if err := processor.WriteToRemote(0, buffer); err != nil {
				log.Error("WriteToRemote error", err)
			}
		}
	}()
	<-common.WaitSignal()
}
