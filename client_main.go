//go:build client

package main

import (
	"flag"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "服务器地址")
	flag.IntVar(&port, "port", 2023, "服务器端口")
}

func main() {

}
