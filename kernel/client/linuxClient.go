package client

import (
	"fmt"
	"github.com/songgao/water/waterutil"
	"net"
	"tiangong/kernel/model"
	"tiangong/kernel/tun"
)

func StartLinuxClientOfUdp(config *model.Config) {
	//创建tun
	myTun := tun.CreateLinuxTun(config.CIDR)

	//监听本地端口，使用UDP
	localAddr, err := net.ResolveUDPAddr("udp", config.LocalAddr)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	defer myTun.Close()
	//读取发送到本地端口报文并写入tun
	go func() {
		buf := make([]byte, 1500)
		for {
			size, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if size == 0 {
				continue
			}
			packet := buf[:size]
			if !waterutil.IsIPv4(packet) {
				continue
			}
			myTun.Write(packet)
		}
	}()
	//读取tun数据并发送到服务端
	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 1500)
	for {
		size, err := myTun.Read(buf)
		if err != nil || size == 0 {
			continue
		}
		packet := buf[:size]
		if !waterutil.IsIPv4(packet) {
			continue
		}
		conn.WriteToUDP(packet, serverAddr)
	}
}
