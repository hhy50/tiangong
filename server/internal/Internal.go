package internal

import (
	"math"
	"tiangong/common/net"
)

var (
	Increment = Incrementer{
		Range: Range{0, math.MaxUint8},
	}
)

func GeneraInternalIp() net.IpAddress {
	return net.IpAddress{
		192, 168, 31, byte(Increment.Next()),
	}
}
