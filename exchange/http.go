package exchange

import (
	"context"
	"github.com/advanced-go/core/controller"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
)

func Http(ctx context.Context, method, url string, h http.Header, body io.Reader) (*http.Response, *runtime.Status) {
	req, _ := http.NewRequestWithContext(ctx, method, url, body)
	if h != nil {
		req.Header = h
	}
	ctrl, status1 := controller.Lookup(url)
	if status1.OK() {
		return ctrl.Do(Do, req)
	}
	return Do(req)
}
