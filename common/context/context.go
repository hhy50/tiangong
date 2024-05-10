package context

import (
	"context"
	"time"
)

var (
	EmptyCtx = context.Background()
)

type Context interface {
	context.Context
	AddValue(name any, value any)
	Cancel()
}

type CancelContext struct {
	context.Context
	values map[any]any
	cancel context.CancelFunc
}

func (c *CancelContext) AddValue(name any, value any) {
	c.values[name] = value
}

func (c *CancelContext) Value(name any) any {
	if v, f := c.values[name]; f {
		return v
	}
	return c.Context.Value(name)
}

func (c *CancelContext) Cancel() {
	c.cancel()
}

func Empty() Context {
	return WithParent(EmptyCtx)
}

func WithParent(parent context.Context) Context {
	ctx, cancel := context.WithCancel(parent)
	return &CancelContext{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}

func WithTimeout(parent context.Context, duration time.Duration) Context {
	ctx, cancel := context.WithTimeout(parent, duration)
	return &CancelContext{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}
