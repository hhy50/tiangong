package context

import (
	"context"
	"time"
)

var (
	EmptyCtx = context.Background()
)

type (
	CancelFunc = context.CancelFunc
)

type Context struct {
	context.Context
	values map[any]any
	cancel context.CancelFunc
}

func (c *Context) AddValue(name any, value any) {
	// c.Context = context.WithValue(c.Context, name, value)
	c.values[name] = value
}

func (c *Context) Value(name any) any {
	if v, f := c.values[name]; f {
		return v
	}
	if p, ok := c.Context.(*Context); ok {
		if v := p.Value(name); v != nil {
			return v
		}
	}
	return c.Context.Value(name)
}

func (c *Context) Cancel() {
	c.cancel()
}

func Empty() Context {
	return WithParent(EmptyCtx)
}

func WithParent(parent context.Context) Context {
	ctx, cancel := context.WithCancel(parent)
	return Context{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}

func WithTimeout(parent Context, duration time.Duration) Context {
	ctx, cancel := context.WithTimeout(&parent, duration)
	return Context{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}
