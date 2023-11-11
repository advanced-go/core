package http2

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

// DoHandler - function type for a Do handler
type DoHandler func(ctx any, r *http.Request, body any) (any, *runtime.Status)

func DoHandlerProxy(ctx any) func(ctx any, r *http.Request, body any) (any, *runtime.Status) {
	switch ptr := ctx.(type) {
	case context.Context:
		if proxies, ok := runtime.IsProxyable(ptr); ok {
			do := findDoProxy(proxies)
			if do != nil {
				return do
			}
		}
	case *http.Request:
		if proxies, ok := runtime.IsProxyable(ptr.Context()); ok {
			do := findDoProxy(proxies)
			if do != nil {
				return do
			}
		}
	}
	return nil
}

func findDoProxy(proxies []any) func(ctx any, r *http.Request, body any) (any, *runtime.Status) {
	for _, p := range proxies {
		if fn, ok := p.(func(ctx any, r *http.Request, body any) (any, *runtime.Status)); ok {
			return fn
		}
	}
	return nil
}
