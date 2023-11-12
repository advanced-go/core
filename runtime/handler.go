package runtime

import (
	"net/http"
)

// HttpHandler - function type for HTTP handling
type HttpHandler func(ctx any, w http.ResponseWriter, r *http.Request) *Status

// DoHandler - function type for a Do handler
type DoHandler func(ctx any, r *http.Request, body any) (any, *Status)
