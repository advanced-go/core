package startup

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// AccessLogFn - typedef for a function that provides access logging
type AccessLogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string)

// ContentMap - slice of any content to be included in a message
type ContentMap map[string][]any

// Credentials - type for a credentials function
type Credentials func() (username string, password string, err error)

// ControllerApply - type for applying a controller
type ControllerApply func(ctx context.Context, statusCode func() int, uri, requestId, method string) (fn func(), newCtx context.Context, rateLimited bool)

// Resource - struct for a resource
type Resource struct {
	Uri string
}

// AccessCredentials - access function for Credentials in a message
func AccessCredentials(msg *Message) Credentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(Credentials); ok {
			return fn
		}
	}
	return nil
}

// AccessControllerApply - access function for ControllerApply in a message
func AccessControllerApply(msg *Message) ControllerApply {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(ControllerApply); ok {
			return fn
		}
	}
	return nil
}

// AccessResource - access function for a resource in a message
func AccessResource(msg *Message) Resource {
	if msg == nil || msg.Content == nil {
		return Resource{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(Resource); ok {
			return url
		}
	}
	return Resource{}
}

func NewStatusCode(status **runtime.Status) func() int {
	return func() int {
		return int((*(status)).Code())
	}
}
