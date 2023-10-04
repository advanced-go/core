package runtime

import (
	"context"
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
