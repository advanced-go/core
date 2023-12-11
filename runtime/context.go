package runtime

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	XRequestId      = "x-request-id"
	ContentLocation = "Content-Location"
	FileScheme      = "file://"
)

type requestContextKey struct{}
type contentLocationT struct{}
type fileUrlLocationT struct{}

var (
	requestKey         = requestContextKey{}
	contentLocationKey = contentLocationT{}
	fileUrlKey         = fileUrlLocationT{}
)

// NewRequestIdContext - creates a new Context with a request id
func NewRequestIdContext(ctx context.Context, requestId string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	} else {
		i := ctx.Value(requestKey)
		if i != nil {
			return ctx
		}
	}
	if requestId == "" {
		requestId = uuid.New().String()
	}
	return ContextWithValue(ctx, requestKey, requestId)
}

// NewRequestContext - creates a new Context with a request id from the request headers
func NewRequestContext(req *http.Request) context.Context {
	if req == nil || req.Header == nil {
		return context.Background()
	}
	if req.Header.Get(XRequestId) == "" {
		req.Header.Add(XRequestId, uuid.New().String())
	}
	return NewRequestIdContext(req.Context(), req.Header.Get(XRequestId))
}

// RequestIdFromContext - return the requestId from a context
func RequestIdFromContext(ctx any) string {
	if ctx == nil {
		return ""
	}
	if ctx2, ok := ctx.(context.Context); ok {
		i := ctx2.Value(requestKey)
		if requestId, ok2 := i.(string); ok2 {
			return requestId
		}
	}
	return ""
}

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

// RequestId - return a request id from any type
func RequestId(t any) string {
	if t == nil {
		return ""
	}
	switch ptr := t.(type) {
	case string:
		return ptr
	case context.Context:
		return RequestIdFromContext(ptr)
	case *http.Request:
		return ptr.Header.Get(XRequestId)
	case http.Header:
		return ptr.Get(XRequestId)
	case Status:
		return ptr.RequestId()
	}
	return ""
}

// GetOrCreateRequestId - return a request id from any type, creating a new id if needed
func GetOrCreateRequestId(t any) string {
	requestId := RequestId(t)
	if requestId == "" {
		s, _ := uuid.NewUUID()
		requestId = s.String()
	}
	return requestId
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
	if location, ok := i.(string); ok {
		if strings.HasPrefix(location, FileScheme) {
			return location, true
		}
	}
	return "", false
}

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
