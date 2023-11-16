package http2

import (
	"context"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

// DoHandler - function type for a Do handler
//type DoHandler func(ctx any, r *http.Request, body any) (any, *runtime.Status)

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

/*
// HttpHandlerProxy - function for finding an HTTP handler proxy
func HttpHandlerProxy(ctx context.Context) func(w http.ResponseWriter, r *http.Request) *runtime.Status {
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		p := findHttpProxy(proxies)
		if p != nil {
			return p
		}
	}

		switch ptr := ctx.(type) {
		case context.Context:
			if proxies, ok := runtime.IsProxyable(ptr); ok {
				p := findHttpProxy(proxies)
				if p != nil {
					return p
				}
			}
		case *http.Request:
			if proxies, ok := runtime.IsProxyable(ptr.Context()); ok {
				p := findHttpProxy(proxies)
				if p != nil {
					return p
				}
			}
		}
	return nil
}

func findHttpProxy(proxies []any) func(w http.ResponseWriter, r *http.Request) *runtime.Status {
	for _, p := range proxies {
		if fn, ok := p.(func(w http.ResponseWriter, r *http.Request) *runtime.Status); ok {
			return fn
		}
	}
	return nil
}

// PostHandlerProxy - function for finding a Post handler proxy
func PostHandlerProxy(ctx context.Context) func(r *http.Request, body any) (any, *runtime.Status) {
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		do := findPostProxy(proxies)
		if do != nil {
			return do
		}
	}

	switch ptr := ctx.(type) {
		case context.Context:
			if proxies, ok := runtime.IsProxyable(ptr); ok {
				do := findPostProxy(proxies)
				if do != nil {
					return do
				}
			}
		case *http.Request:
			if proxies, ok := runtime.IsProxyable(ptr.Context()); ok {
				do := findPostProxy(proxies)
				if do != nil {
					return do
				}
			}
		}

	return nil
}

func findPostProxy(proxies []any) func(r *http.Request, body any) (any, *runtime.Status) {
	for _, p := range proxies {
		if fn, ok := p.(func(r *http.Request, body any) (any, *runtime.Status)); ok {
			return fn
		}
	}
	return nil
}

// GetHandler - function type for a Get handler
//type GetHandler func(ctx any, uri, variant string) (any, *runtime.Status)

// GetHandlerProxy - function for finding a Get handler proxy
func GetHandlerProxy(ctx context.Context) func(h http.Header, uri string) (any, *runtime.Status) {
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		do := findGetProxy(proxies)
		if do != nil {
			return do
		}
	}
			switch ptr := ctx.(type) {
		case context.Context:
			if proxies, ok := runtime.IsProxyable(ptr); ok {
				do := findGetProxy(proxies)
				if do != nil {
					return do
				}
			}
		case *http.Request:
			if proxies, ok := runtime.IsProxyable(ptr.Context()); ok {
				do := findGetProxy(proxies)
				if do != nil {
					return do
				}
			}
		}

	return nil
}

func findGetProxy(proxies []any) func(h http.Header, uri string) (any, *runtime.Status) {
	for _, p := range proxies {
		if fn, ok := p.(func(h http.Header, uri string) (any, *runtime.Status)); ok {
			return fn
		}
	}
	return nil
}


*/
