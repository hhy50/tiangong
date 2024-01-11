package net_test

import (
	"testing"
	"tiangong/common/net"
)

func TestIpAddress(t *testing.T) {
	if net.Local.String() != "127.0.0.1" {
		t.Error("net.Local.String() != \"127.0.0.1\"")
		return
	}

	address := net.IpAddress{192, 168, 1, 1}
	if address.String() != "192.168.1.1" {
		t.Error("net.Local.String() != \"192.168.1.1\"")
		return
	}
}
