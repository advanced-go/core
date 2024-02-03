package runtime

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	XRequestId = "x-request-id"
	XRelatesTo = "x-relates-to"
)

type requestContextKey struct{}

var (
	requestKey = requestContextKey{}
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
	return context.WithValue(ctx, requestKey, requestId)
}

// NewRequestIdContextFromHeader - creates a new Context with a request id from the request headers
func NewRequestIdContextFromHeader(h http.Header) context.Context {
	if h == nil {
		return context.Background()
	}
	if h.Get(XRequestId) == "" {
		h.Add(XRequestId, uuid.New().String())
	}
	return NewRequestIdContext(context.Background(), h.Get(XRequestId))
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
	}
	return ""
}

// GetOrCreateRequestId2 - return a request id from any type, creating a new id if needed
func GetOrCreateRequestId2(t any) string {
	requestId := RequestId(t)
	if requestId == "" {
		s, _ := uuid.NewUUID()
		requestId = s.String()
	}
	return requestId
}

// AddRequestId2 - add a request to an http.Request or an http.Header
func AddRequestId2(t any) http.Header {
	if req, ok := t.(*http.Request); ok {
		req.Header = addRequestId(req.Header)
		return req.Header
	}
	if h, ok := t.(http.Header); ok {
		return addRequestId(h)
	}
	return make(http.Header)
}

func addRequestId(h http.Header) http.Header {
	if h == nil {
		h = make(http.Header)
	}
	id := h.Get(XRequestId)
	if len(id) == 0 {
		uid, _ := uuid.NewUUID()
		id = uid.String()
		h.Set(XRequestId, id)
	}
	return h
}
