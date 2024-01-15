package tun

import (
	"fmt"
	"github.com/songgao/water"
	"os"
	"os/exec"
)

func CreateLinuxTun(cidr string) (myTun *water.Interface) {
	tunName := "my-tun"

	config := water.Config{
		DeviceType: water.TUN,
	}

	myTun, err := water.New(config)

	if err != nil {
		fmt.Println(err)
	}

	execCmd("/sbin/ip", "link", "set", "dev", tunName, "mtu", "1500")
	execCmd("/sbin/ip", "addr", "add", cidr, "dev", tunName)
	execCmd("/sbin/ip", "link", "set", "dev", tunName, "up")

	return myTun
}

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
