package net

import "net"

var Dial = net.Dial

type Conn struct {
	net.Conn
}
