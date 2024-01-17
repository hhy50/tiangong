package client

import "tiangong/common/net"

type Config struct {
	ServerHost string
	ServerPort net.Port
	Export     string
}

func defaultValue(key string) string {
	return ""
}
