package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tiangong/common"
	"tiangong/common/errors"
	"tiangong/common/log"
	_ "tiangong/common/log"
	"tiangong/kernel"
	"tiangong/server"
	"time"
)

var (
	cp string
)

func init() {
	banner :=`
	_____ _                    ____                     ____
	|_   _(_) __ _ _ __        / ___| ___  _ __   __ _  / ___|  ___ _ ____   _____ _ __
	  | | | |/ _| | '_ \ _____| |  _ / _ \| '_ \ / _| | \___ \ / _ \ '__\ \ / / _ \ '__|
	  | | | | (_| | | | |_____| |_| | (_) | | | | (_| |  ___) |  __/ |   \ V /  __/ |
	  |_| |_|\__,_|_| |_|      \____|\___/|_| |_|\__, | |____/ \___|_|    \_/ \___|_|
	                                             |___/
	Kernel Version: %s Pid:%d Now: %s
`
	fmt.Printf(banner, kernel.Version(), os.Getpid(), time.Now().Format(common.DateFormat))
	flag.StringVar(&cp, "conf", "", "Config file path")
}

func main() {
	flag.Parse()
	log.Debug("start...");
	server, err := server.NewServer(cp)
	if err != nil {
		if e, ok := err.(*errors.Error); ok {
			fmt.Printf("[error] %s", e.Error())
			return
		}
		panic(err)
	}
	server.Start()
	defer server.Stop()

	pauseProcesss()
}

func pauseProcesss() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
}
