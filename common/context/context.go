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

type EnhanceContext struct {
	context.Context
	values map[any]any
	cancel context.CancelFunc
}

func (c *EnhanceContext) AddValue(name any, value any) {
	c.values[name] = value
}

func (c *EnhanceContext) Value(name any) any {
	if v, f := c.values[name]; f {
		return v
	}
	return c.Context.Value(name)
}

//
//func Done() <-chan struct{} {
//
//}

func (c *EnhanceContext) Cancel() {
	c.cancel()
}

func Empty() Context {
	return WithParent(EmptyCtx)
}

func WithParent(parent context.Context) Context {
	ctx, cancel := context.WithCancel(parent)
	return &EnhanceContext{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}

func WithTimeout(parent context.Context, duration time.Duration) Context {
	ctx, cancel := context.WithTimeout(parent, duration)
	return &EnhanceContext{
		Context: ctx,
		cancel:  cancel,
		values:  map[any]any{},
	}
}
