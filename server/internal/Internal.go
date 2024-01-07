package internal

import "tiangong/common/net"

var (
	Increment = Incrementer{
		Range: Range{0, 255},
	}
)

func GeneraInternalIp() net.IpAddress {
	return net.IpAddress{
		192, 168, 31, byte(Increment.Next()),
	}
}
