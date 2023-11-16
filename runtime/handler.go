package runtime

import (
	"context"
	"net/http"
)

// HttpHandler - function type for HTTP handling
type HttpHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) *Status

// DoHandler - function type for a Do handler
type DoHandler func(ctx any, r *http.Request, body any) (any, *Status)

// PostHandler - function type for a Post handler
type PostHandler func(ctx context.Context, r *http.Request, body any) (any, *Status)
