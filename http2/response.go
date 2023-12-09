package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func NewResponse(status runtime.Status) *http.Response {
	r := new(http.Response)
	r.StatusCode = status.Http()
	r.Status = status.Description()
	return r
}
