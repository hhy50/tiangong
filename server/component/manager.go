package component

import (
	"context"

	"github.com/haiyanghan/tiangong/common/errors"
)

type CreatorFunc = func(context.Context) (Component, error)

var (
	components = map[string]CreatorFunc{}
)

func Register(name string, creator CreatorFunc) error {
	if _, f := components[name]; f {
		return errors.NewError("Duplicate component "+name, nil)
	}
	components[name] = creator
	return nil
}

func GetComponents() map[string]CreatorFunc {
	return components
}
