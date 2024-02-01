package main

import (
	"flag"

	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/kernel/proxy"
)

var (
	cp string
)

func init() {
	flag.StringVar(&cp, "conf", "", "config file path")
}

func main() {
	flag.Parse()
	log.InitLog()

	processor := NewProcessor(cp)
	if err := processor.Start(); err != nil {
		log.Error("Kernel client start error, ", err)
		return
	}

	defer proxy.ResetProxy()
	log.Info("Kernel client start success")
	<-common.WaitSignal()
}
