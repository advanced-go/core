package runtime

import (
	"context"
	"net/http"
	"time"
)

type proxyContext struct {
	ctx     context.Context
	proxies []any
}

func (p *proxyContext) Deadline() (deadline time.Time, ok bool) {
	return p.ctx.Deadline()
}

func (p *proxyContext) Done() <-chan struct{} {
	return p.ctx.Done()
}

func (p *proxyContext) Err() error {
	return p.ctx.Err()
}

func (p *proxyContext) Value(key any) any {
	return p.ctx.Value(key)
}

func (p *proxyContext) getProxies() []any {
	return p.proxies
}

func (p *proxyContext) add(proxy any) {
	if proxy != nil {
		p.proxies = append(p.proxies, proxy)
	}
}

func (p *proxyContext) withValue(key, val any) context.Context {
	p.ctx = context.WithValue(p.ctx, key, val)
	return p
}

func DoHandlerProxy(ctx context.Context) DoHandler {
	if proxies, ok := IsProxyable(ctx); ok {
		do := findDoProxy(proxies)
		if do != nil {
			return do
		}
	}
	return nil
}

func findDoProxy(proxies []any) func(ctx any, r *http.Request, body any) (any, *Status) {
	for _, p := range proxies {
		if fn, ok := p.(func(ctx any, r *http.Request, body any) (any, *Status)); ok {
			return fn
		}
	}
	return nil
}
