//go:build server

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/haiyanghan/tiangong"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
	"github.com/haiyanghan/tiangong/common/log"
	"github.com/haiyanghan/tiangong/server"
)

var (
	cp string
)

func init() {
	banner := `
	 _____ _                    ____                     ____
	|_   _(_) __ _ _ __        / ___| ___  _ __   __ _  / ___|  ___ _ ____   _____ _ __
	  | | | |/ _| | '_ \ _____| |  _ / _ \| '_ \ / _| | \___ \ / _ \ '__\ \ / / _ \ '__|
	  | | | | (_| | | | |_____| |_| | (_) | | | | (_| |  ___) |  __/ |   \ V /  __/ |
	  |_| |_|\__,_|_| |_|      \____|\___/|_| |_|\__, | |____/ \___|_|    \_/ \___|_|
	                                             |___/
	TianGong Version: %s Pid:%d Now: %s
`
	fmt.Printf(banner, tiangong.Version(), os.Getpid(), time.Now().Format(common.DateFormat))
	flag.StringVar(&cp, "conf", "", "Config file path")
}

func main() {
	flag.Parse()
	log.InitLog()

	server, err := server.NewServer(cp)
	if err != nil {
		log.Error("Init TianGong Server error, ", err)
		handlerError(err)
		return
	}

	if err = server.Start(); err != nil {
		handlerError(err)
		return
	}

	defer server.Stop()
	log.Info("TianGong Server started")
	<-common.WaitSignal()
}

func handlerError(err error) {
	if e, ok := err.(*errors.Error); ok {
		log.Error("TianGong Server start error, ", e)
		return
	}
	panic(err)
}
