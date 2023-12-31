package context

import "context"

var (
	empty = Context{Context: context.Background()}
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

func Background() Context {
	return empty
}

func Wrap(c context.Context) Context {
	return Context{Context: c}
}

func WithCancel(p Context) (Context, CancelFunc) {
	ctx, cf := context.WithCancel(p)
	return Wrap(ctx), cf
}
