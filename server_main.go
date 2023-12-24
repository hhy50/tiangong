package main

import (
	"flag"
	"fmt"
	"tiangong/common"
	"tiangong/kernel"
	"tiangong/server"
	"time"
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
	Kernel Version: %s, Now: %s
`
	fmt.Printf(banner, kernel.Version(), time.Now().Format(common.DateFormat))
	flag.StringVar(&cp, "conf", "", "config file path")
}

func main() {
	server, err := server.NewServer(cp)
	if err != nil {
		panic(err)
	}
	server.Start()
	defer server.Stop()
}
