//go:build client

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/haiyanghan/tiangong"
	"github.com/haiyanghan/tiangong/client"
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/conf"
	"github.com/haiyanghan/tiangong/common/log"
)

var ()

func init() {
	banner := `
	  _____ _                    ____                      ____ _     ___ _____ _   _ _____
	 |_   _(_) __ _ _ __        / ___| ___  _ __   __ _   / ___| |   |_ _| ____| \ | |_   _|
	   | | | |/ _| | '_ \ _____| |  _ / _ \| |_ \ / _| | | |   | |    | ||  _| |  \| | | |
	   | | | | (_| | | | |_____| |_| | (_) | | | | (_| | | |___| |___ | || |___| |\  | | |
	   |_| |_|\__,_|_| |_|      \____|\___/|_| |_|\__, |  \____|_____|___|_____|_| \_| |_|
	                                               |___/
		TianGong Version: %s Pid:%d Now: %s
`
	fmt.Printf(banner, tiangong.Version(), os.Getpid(), time.Now().Format(common.DateFormat))
}

func main() {
	flag.Parse()
	log.InitLog()
	conf.Load()

	log.Info("TianGong Client start...")
	c, err := client.NewClient()
	if err != nil {
		log.Error("Create new client error", err)
		return
	}
	if err := c.Start(); err != nil {
		log.Error("Client start error, ", err)
		return
	}
	log.Info("TianGong Client started")
	<-common.WaitSignal()
}
