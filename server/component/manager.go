package component

import (
	"github.com/haiyanghan/tiangong/common/context"
	"github.com/haiyanghan/tiangong/common/errors"
)

type CreatorFunc = func(context.Context) (Component, error)

var (
	components = map[string]CreatorFunc{}
)

func Register(name string, creator CreatorFunc) {
	if _, f := components[name]; f {
		panic(errors.NewError("Duplicate component "+name, nil))
	}
	components[name] = creator
}

func GetComponents() map[string]CreatorFunc {
	return components
}
