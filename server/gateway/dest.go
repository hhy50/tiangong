package gateway

import "tiangong/common/net"

type Destination interface {
}

type DirectDestination struct {
	IpAddress string
	Port      net.Port
}
