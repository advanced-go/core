package http2test

import (
	"github.com/advanced-go/core/exchange"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

// DoT - do a Http exchange with deserialization
func DoT[T any](req *http.Request) (resp *http.Response, t T, status runtime.Status) {
	resp, status = exchange.Do(req)
	if !status.OK() {
		return nil, t, status
	}
	//t, status = http2.Deserialize[T](resp.Body)
	return
}
