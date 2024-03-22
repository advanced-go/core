package exchange

import (
	"context"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	getLocation = PkgPath + ":Get"
)

// Get - process an HTTP Get request
func Get(ctx context.Context, uri string, h http.Header) (resp *http.Response, status *runtime.Status) {
	if len(uri) == 0 {
		return serverErrorResponse(), runtime.NewStatusError(http.StatusBadRequest, errors.New("error: URI is empty"))
	}
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return serverErrorResponse(), runtime.NewStatusError(http.StatusBadRequest, err)
	}
	if h != nil {
		req.Header = h
	}
	// exchange.Do() will always return a non nil *http.Response
	resp, status = Do(req)
	if !status.OK() {
		status.AddLocation()
	}
	return
}
