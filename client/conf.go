package client

import (
	"github.com/haiyanghan/tiangong/common"
	"github.com/haiyanghan/tiangong/common/errors"
)

type Config struct {
	Address  string `prop:"server.address"`
	Key      string `prop:"server.key"`
	Internal string `prop:"server.internal"`
	Name     string `prop:"main.name"`
	Export   string `prop:"main.export"`
}

func (c Config) Require() error {
	if common.IsEmpty(c.Address) {
		return errors.NewError("server.address not be null", nil)
	}

	if common.IsEmpty(c.Key) {
		return errors.NewError("server.key not be null", nil)
	}
	return nil
}

func defaultValue(key string) string {
	return ""
}
