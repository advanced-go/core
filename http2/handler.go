package http2

import (
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

// TypeHandler - function type for TypeHandler handler
type TypeHandler func(r *http.Request, body any) (any, *runtime.Status)

// DoHandler - function type for a Do handler
type DoHandler func(ctx any, r *http.Request, body any) (any, *runtime.Status)
