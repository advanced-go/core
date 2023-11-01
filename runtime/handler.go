package runtime

import "net/http"

// TypeHandlerFn - function type for TypeHandler connector
type TypeHandlerFn func(r *http.Request, body any) (any, *Status)
