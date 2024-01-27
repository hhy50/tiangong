package server

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/log"
	"strconv"

	"github.com/google/uuid"
)

var (
	HostDef = common.Pair[string, string]{
		First:  "host",
		Second: "0.0.0.0",
	}

	SrvPortDef = common.Pair[string, int]{
		First:  "srvPort",
		Second: 2023,
	}

	HttpPortDef = common.Pair[string, int]{
		First:  "httpPort",
		Second: 2024,
	}

	UserNameDef = common.Pair[string, string]{
		First:  "userName",
		Second: "admin",
	}

	PasswdDef = common.Pair[string, string]{
		First:  "passwd",
		Second: "",
	}
)

type Config struct {
	Host     string `prop:"host"`
	SrvPort  int    `prop:"srvPort"`
	HttpPort int    `prop:"httpPort"`
	UserName string `prop:"username"`
	Passwd   string `prop:"passwd"`
	Key      string `prop:"key"`
}

// defaultValue
func defaultValue(key string) string {
	switch key {
	case HostDef.First:
		return HostDef.Second
	case SrvPortDef.First:
		return strconv.Itoa(SrvPortDef.Second)
	case HttpPortDef.First:
		return strconv.Itoa(HttpPortDef.Second)
	case UserNameDef.First:
		return UserNameDef.Second
	case PasswdDef.First:
		passwd := uuid.New().String()
		log.Warn("httpPasswd is not set, Generate a random password: %s", passwd)
		return passwd
	default:
		return ""
	}
}
