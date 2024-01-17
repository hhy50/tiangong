//go:build client

package main

import (
	"flag"
	"os"
	"tiangong/client"
	"tiangong/common"
	"tiangong/common/log"
)

var (
	cp string
)

func init() {
	banner := `
	  _____ _                    ____                      ____ _     ___ _____ _   _ _____ 
	 |_   _(_) __ _ _ __        / ___| ___  _ __   __ _   / ___| |   |_ _| ____| \ | |_   _|
	   | | | |/ _| | '_ \ _____| |  _ / _ \| |_ \ / _| | | |   | |    | ||  _| |  \| | | |  
	   | | | | (_| | | | |_____| |_| | (_) | | | | (_| | | |___| |___ | || |___| |\  | | |  
	   |_| |_|\__,_|_| |_|      \____|\___/|_| |_|\__, |  \____|_____|___|_____|_| \_| |_|  
												  |___/                                     
		Kernel Version: %s Pid:%d Now: %s
`
	fmt.Printf(banner, kernel.Version(), os.Getpid(), time.Now().Format(common.DateFormat))
	flag.StringVar(&cp, "conf", "", "Config file path")
}

func main() {
	flag.Parse()
	log.InitLog()
	log.Info("TianGong Client start...")

	client.NewClient(cp)

	log.Info("TianGong Client started")
	<-common.WaitSignal()
}
