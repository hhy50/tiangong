package server

import (
	"fmt"
	"github.com/songgao/water/waterutil"
	"net"
	"tiangong/kernel/model"
	"tiangong/kernel/tun"
	netutil "tiangong/kernel/util"
	"time"

	"github.com/patrickmn/go-cache"
)

func StartLinuxServerOfUdp(config *model.Config) {
	//根据路由缓存连接
	clientCache := cache.New(10*time.Minute, 15*time.Minute)
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
	//从本地tun中读取数据，并发回client
	go func() {
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
			srcAddr, dstAddr := netutil.GetAddr(packet)
			if srcAddr == "" || dstAddr == "" {
				continue
			}
			key := fmt.Sprintf("%v->%v", dstAddr, srcAddr)
			clientAddr, ok := clientCache.Get(key)
			if ok {
				conn.WriteToUDP(packet, clientAddr.(*net.UDPAddr))
			}
		}
	}()
	//读取client发送来的数据，并写入本地tun
	buf := make([]byte, 1500)
	for {
		size, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil || size == 0 {
			continue
		}
		packet := buf[:size]
		if !waterutil.IsIPv4(packet) {
			continue
		}
		myTun.Write(packet)
		//从ip报文中获取源IP和目标IP
		srcAddr, dstAddr := netutil.GetAddr(packet)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		//根据IP地址缓存
		key := fmt.Sprintf("%v->%v", srcAddr, dstAddr)
		clientCache.Set(key, clientAddr, cache.DefaultExpiration)
	}
}
