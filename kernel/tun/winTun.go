package tun

import (
	"fmt"
	"golang.zx2c4.com/wireguard/tun"
	"net/netip"
	"tiangong/kernel/winipcfg"
)

func CreateWinTun(cidr string) (myTun tun.Device) {
	tunName := "my-tun"

	myTun, err := tun.CreateTUN(tunName, 0)

	if err != nil {
		fmt.Println(err)
	}

	ip, err := netip.ParsePrefix(cidr)
	if err != nil {
		fmt.Println(err)
	}

	link := winipcfg.LUID(myTun.(*tun.NativeTun).LUID())
	err = link.SetIPAddresses([]netip.Prefix{ip})
	if err != nil {
		fmt.Println(err)
	}

	return myTun
}
