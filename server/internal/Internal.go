package internal

import (
	"math"
	"tiangong/common"
	"tiangong/common/net"
)

var (
	Increment = common.Incrementer{
		Range: common.Range{0, math.MaxUint8},
	}
)

func GeneraInternalIp() net.IpAddress {
	return net.IpAddress{
		192, 168, 31, byte(Increment.Next()),
	}
}
