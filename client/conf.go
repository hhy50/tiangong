package client

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
)

type Config struct {
	ServerHost string `prop:"serverHost"`
	ServerPort int    `prop:"serverPort"`
	Key        string `prop:"key"`
	Name       string `prop:"name"`
	Internal   string `prop:"internal"`
	Export     string `prop:"export"`
}

func (c Config) Require() error {
	if common.IsEmpty(c.ServerHost) {
		return errors.NewError("serverHost not be null", nil)
	}

	if c.ServerPort == 0 {
		return errors.NewError("serverPort not be null", nil)
	}

	if common.IsEmpty(c.Key) {
		return errors.NewError("server Key not be null", nil)
	}
	return nil
}

func defaultValue(key string) string {
	return ""
}
