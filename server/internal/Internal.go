package internal

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/net"
	"math"
)

var (
	Increment = common.Incrementer{
		Range: common.Range{0, math.MaxUint8},
	}
)

func GeneraInternalIp() net.IpAddress {
	return net.IpAddress{
		172, 1, 0, byte(Increment.Next()),
	}
}
