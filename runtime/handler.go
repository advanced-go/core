package runtime

import (
	"net/http"
)

// TypeHandlerFn - function type for TypeHandler connector
type TypeHandlerFn func(r *http.Request, body any) (any, *Status)

// DoHandlerFn - function type for a Do handler
type DoHandlerFn func(ctx any, r *http.Request, body any) (any, *Status)
