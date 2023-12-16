package runtimetest

import (
	"context"
	"net/http"
	"strings"
)

const (
	FileScheme      = "file://"
	UrnScheme       = "urn:"
	ContentLocation = "Content-Location"
)

type contentLocationT struct{}
type fileUrlLocationT struct{}

var (
	contentLocationKey = contentLocationT{}
	fileUrlKey         = fileUrlLocationT{}
)

// NewProxyContext - create a new Context interface, containing a proxy
func NewProxyContext(ctx context.Context, proxy any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	} else {
		if pCtx, ok := any(ctx).(*proxyContext); ok {
			pCtx.add(proxy)
			return ctx
		}
	}
	mux := new(proxyContext)
	mux.ctx = ctx
	if proxy != nil {
		mux.proxies = append(mux.proxies, proxy)
	}
	return mux
}

// ContextWithValue - create a new context with a value, updating the context if it is a Proxy context
func ContextWithValue(ctx context.Context, key any, val any) context.Context {
	if ctx == nil {
		return nil
	}
	if pCtx, ok := any(ctx).(*proxyContext); ok {
		return pCtx.withValue(key, val)
	}
	return context.WithValue(ctx, key, val)
}

// IsProxyable - determine if the context is a ProxyContext, and return proxies
func IsProxyable(ctx context.Context) ([]any, bool) {
	if ctx == nil {
		return nil, false
	}
	if pCtx, ok := any(ctx).(*proxyContext); ok {
		return pCtx.getProxies(), true
	}
	return nil, false
}

// NewContentLocationContext - creates a new Context with a content location
func NewContentLocationContext(ctx context.Context, h http.Header) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if h == nil {
		return ctx
	}
	v := h.Get(ContentLocation)
	if len(v) == 0 {
		return ctx
	}
	return context.WithValue(ctx, contentLocationKey, v)
}

// ContentLocationFromContext - return a content location from a context
func ContentLocationFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	i := ctx.Value(contentLocationKey)
	if i == nil {
		return "", false
	}
	if location, ok := i.(string); ok {
		if strings.HasPrefix(location, FileScheme) {
			return location, true
		}
	}
	return "", false
}

// NewFileUrlContext - creates a new Context with a content location
func NewFileUrlContext(ctx context.Context, url string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if len(url) == 0 {
		return ctx
	}
	return context.WithValue(ctx, fileUrlKey, url)
}

// FileUrlFromContext - return a content location from a context
func FileUrlFromContext(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	i := ctx.Value(fileUrlKey)
	if i == nil {
		return "", false
	}
	if uri, ok := i.(string); ok {
		if strings.HasPrefix(uri, FileScheme) || strings.HasPrefix(uri, UrnScheme) {
			return uri, true
		}
	}
	return "", false
}

/*
type statusT struct{}

var (
	statusKey = statusT{}
)

// NewStatusContext - creates a new Context with a Status
func NewStatusContext(ctx context.Context, status Status) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, statusKey, status)
}

// StatusFromContext - return a Status from a context2
func StatusFromContext(ctx context.Context) Status {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(statusKey)
	if i == nil {
		return nil
	}
	if status, ok := i.(Status); ok {
		return status
	}
	return nil
}




*/
