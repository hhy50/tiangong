package common

import "context"

var (
	EmptyCtx = context.Background()

	ListenerCtxKey = "listener"
)

type CancelFunc = context.CancelFunc

type Context struct {
	context.Context
}

func (c *Context) Add(key interface{}, value interface{}) {
	*c = Context{
		Context: context.WithValue(c.Context, key, value),
	}
}

func Wrap(c context.Context) Context {
	return Context{Context: c}
}
