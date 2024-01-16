package client

import (
	"fmt"
	"github.com/songgao/water/waterutil"
	"net"
	"tiangong/kernel/model"
	"tiangong/kernel/tun"
)

func StartWinClient(config *model.Config) {
	myTun := tun.CreateWinTun(config.CIDR)

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
			packets := make([][]byte, 1)
			packets[0] = packet
			myTun.Write(packets, 0)
		}
	}()

	serverAddr, err := net.ResolveUDPAddr("udp", config.ServerAddr)
	if err != nil {
		fmt.Println(err)
	}
	bufs := make([][]byte, 1)
	bufs[0] = make([]byte, 1500)

	sizes := make([]int, 1)
	sizes[0] = 1500 + 1
	for {
		sizes[0] = 1500 + 1

		result, err := myTun.Read(bufs, sizes, 0)

		if err != nil {
			fmt.Println(err)
			continue
		}
		if result != 1 {
			continue
		}
		packet := bufs[0][:sizes[0]]
		if !waterutil.IsIPv4(packet) {
			continue
		}
		conn.WriteToUDP(packet, serverAddr)
	}
}
