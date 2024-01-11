package main

import (
	"flag"
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
	flag.StringVar(&server, "server", "", "Specify target server")
	flag.IntVar(&port, "port", 2023, "Specify target port")
	flag.StringVar(&token, "token", "", "Token")
	flag.StringVar(&subHost, "subHost", "", "SubHost")
}

func main() {
	log.InitLog()
}
