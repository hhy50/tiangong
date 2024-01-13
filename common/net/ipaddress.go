package net

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IpAddress A four byte IP address
type IpAddress [net.IPv4len]byte

// ConnHandlerFunc connect success exec
type ConnHandlerFunc func(net.Conn) error

var (
	Internal IpAddress
	Local    = IpAddress{127, 0, 0, 1}
)

func init() {
	// TODO
	// Get Internal ip From
}

// Port Range 1024~65535
type Port uint16

func (p IpAddress) GetA() byte {
	return p[0]
}

func (p IpAddress) GetB() byte {
	return p[1]
}

func (p IpAddress) GetC() byte {
	return p[2]
}

func (p IpAddress) GetD() byte {
	return p[3]
}

func (p IpAddress) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", p.GetA(), p.GetB(), p.GetC(), p.GetD())
}

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

func ConvertIp(bytes []byte) IpAddress {
	ip := IpAddress{}
	copy(ip[:], bytes)
	return ip
}

func ParseIp(parse string) IpAddress {
	split := strings.Split(parse, ".")
	if len(split) == net.IPv4len {
		a, _ := strconv.Atoi(split[0])
		b, _ := strconv.Atoi(split[1])
		c, _ := strconv.Atoi(split[2])
		d, _ := strconv.Atoi(split[3])
		return IpAddress{
			byte(a),
			byte(b),
			byte(c),
			byte(d),
		}
	}
	return IpAddress{}
}
