package http2

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

// See https://tools.ietf.org/html/rfc6265 for details of each of the fields of the above cookie.

func ReadCookies(req *http.Request) map[string]*http.Cookie {
	if req == nil {
		return nil
	}
	jar := make(map[string]*http.Cookie)
	for _, c := range req.Cookies() {
		jar[c.Name] = c
	}
	return jar
}

func AddHeaders(req *http.Request, header http.Header) {
	if req == nil || header == nil {
		return
	}
	for key, element := range header {
		req.Header.Add(key, createValue(element))
	}
}

func createValue(v []string) string {
	if v == nil {
		return ""
	}
	var value string
	for i, item := range v {
		if i > 0 {
			value += ","
		}
		value += item
	}
	return value
}

func UpdateHeaders(req *http.Request) *http.Request {
	if req == nil {
		return nil
	}
	AddRequestId(req)
	//if log.AccessFromContext(req.Context()) != nil {
	//	return req
	//}
	//	if fn := log.Access(); fn != nil {
	//		return req.Clone(log.NewAccessContext(req.Context()))
	//	}
	return req
}

func NewRequest(ctx any, method, uri, variant string) (*http.Request, *runtime.Status) {
	newCtx := newContext(ctx)

	// Check for access function
	//if log.AccessFromContext(newCtx) == nil {
	//	if fn := log.Access(); fn != nil {
	//		newCtx = log.NewAccessContext(newCtx)
	//	}
	//}
	// Create request id and add to context
	requestId := newId(ctx)
	if id := runtime.RequestIdFromContext(newCtx); len(id) == 0 {
		newCtx = runtime.NewRequestIdContext(newCtx, requestId)
	}
	req, err := http.NewRequestWithContext(newCtx, method, uri, nil)
	if err != nil {
		return nil, runtime.NewStatusError(http.StatusBadRequest, "/NewRequest", err)
	}
	if len(variant) != 0 {
		req.Header.Add(ContentLocation, variant)
	}
	req.Header.Add(runtime.XRequestId, requestId)
	return req, runtime.NewStatusOK()
}

func newContext(ctx any) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if ctx2, ok := ctx.(context.Context); ok {
		return ctx2
	}
	//if r, ok := ctx.(*http.Request); ok && r.Context() != nil {
	//	return r.Context()
	//}
	return context.Background()
}

func newId(ctx any) string {
	if ctx == nil {
		uid, _ := uuid.NewUUID()
		return uid.String()
	}
	var id = ""
	if r, ok := ctx.(*http.Request); ok {
		id = r.Header.Get(runtime.XRequestId)
		if len(id) == 0 {
			uid, _ := uuid.NewUUID()
			id = uid.String()
		}
		return id
	}
	if ctx2, ok := ctx.(context.Context); ok {
		id = runtime.RequestIdFromContext(ctx2)
		if len(id) == 0 {
			uid, _ := uuid.NewUUID()
			id = uid.String()
		}
	}
	return id
}
