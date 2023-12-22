package main

import (
	"flag"
	"fmt"
	"tiangong/client"
	"time"
)

var (
	host string
	port int
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1", "服务器地址")
	flag.IntVar(&port, "port", 2023, "服务器端口")
}

// go build -o client main.go client.go
// ./client -host localhost -port 2023
func main() {
	flag.Parse()
	tclient := client.NewClient(host, port)
	if tclient == nil {
		fmt.Println("服务器链接失败")
		return
	}
	fmt.Println("服务器链接成功")

	var msg string
	for {
		fmt.Println("请输入要发送的数据")
		_, err := fmt.Scanln(&msg)
		if err != nil {
			fmt.Println("Scanln error", err)
			break
		}
		if msg == "quit" {
			break
		}
		tclient.Write([]byte(msg))
		time.Sleep(1 * time.Second)
	}

	select {}
}
