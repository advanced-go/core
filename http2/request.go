package http2

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
)

// See https://tools.ietf.org/html/rfc6265 for details of each of the fields of the above cookie.

// ReadCookies - read the cookies from a request
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

// NewRequest - create a new Http request adding the request id
// TO DO: if the ctx is a http.Header, then need to add to the new request
func NewRequest(ctx any, method string, uri any, body io.Reader) (*http.Request, *runtime.Status) {
	newCtx := newContext(ctx)

	// Create request id and add to context
	requestId := newId(ctx)
	if id := runtime.RequestIdFromContext(newCtx); len(id) == 0 {
		newCtx = runtime.NewRequestIdContext(newCtx, requestId)
	}
	if len(method) == 0 {
		method = http.MethodGet
	}
	s := "https://somedomain.com/invalid-uri-or-type"
	if url, ok := uri.(*url.URL); ok {
		s = url.String()
	} else {
		if s2, ok2 := uri.(string); ok2 {
			s = s2
		}
	}
	req, err := http.NewRequestWithContext(newCtx, method, s, body)
	if err != nil {
		return nil, runtime.NewStatusError(http.StatusBadRequest, err)
	}
	req.Header.Add(runtime.XRequestId, requestId)
	return req, runtime.StatusOK()
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
	if r, ok := ctx.(http.Header); ok {
		id = r.Get(runtime.XRequestId)
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

// ValidateRequest - validate the request given an embedded URN path
func ValidateRequest(req *http.Request, path string) (string, *runtime.Status) {
	if req == nil {
		return "", runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("error: Request is nil"))
	}
	reqNid, reqPath, ok := uri.UprootUrn(req.URL.Path)
	if !ok {
		return "", runtime.NewStatusError(http.StatusBadRequest, errors.New(fmt.Sprintf("error: invalid URI, path is not valid: \"%v\"", req.URL.Path)))
	}
	if reqNid != path {
		return "", runtime.NewStatusError(http.StatusBadRequest, errors.New(fmt.Sprintf("error: invalid URI, NID does not match: \"%v\" \"%v\"", req.URL.Path, path)))
	}
	return reqPath, runtime.StatusOK()
}
