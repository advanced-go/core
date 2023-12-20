package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

// DoT - do a Http exchange with deserialization
func DoT[T any](req *http.Request) (resp *http.Response, t T, status runtime.Status) {
	resp, status = Do(req)
	if !status.OK() {
		return nil, t, status
	}
	t, status = Deserialize[T](resp.Body)
	return
}
