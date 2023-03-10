package runtime

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type withValue interface {
	withValue(key, val any) context.Context
}

type exchangeContext struct {
	ctx context.Context
	do  func(*http.Request) (*http.Response, error)
}

// ContextWithHttpExchange - create a new Context interface, containing a Http exchange function
func ContextWithHttpExchange(ctx context.Context, do func(*http.Request) (*http.Response, error)) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &exchangeContext{ctx: ctx, do: do}
}

func (c *exchangeContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *exchangeContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *exchangeContext) Err() error {
	return c.ctx.Err()
}

func (c *exchangeContext) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *exchangeContext) Do(req *http.Request) (*http.Response, error) {
	if c.do == nil {
		return nil, errors.New("invalid argument: function Do() is nil")
	}
	return c.do(req)
}

func (c *exchangeContext) withValue(key, val any) context.Context {
	c.ctx = context.WithValue(c.ctx, key, val)
	return c
}
