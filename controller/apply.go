package controller

import (
	"context"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	upstreamTimeoutFlag = "UT"
	//statusDeadlineExceeded = 4
)

func Apply(ctx context.Context, newCtx *context.Context, method, uri, routeName string, h http.Header, duration time.Duration, statusCode access.StatusCodeFunc) func() {
	var cancelFunc context.CancelFunc

	if ctx == nil {
		ctx = context.Background()
	}
	*newCtx = ctx
	start := time.Now()
	req, _ := http.NewRequest(method, uri, nil)
	if h != nil {
		req.Header = h
	}
	// if a timeout and there is no deadline in the current ctx, then create a new context
	if duration > 0 {
		if _, ok := ctx.Deadline(); !ok {
			*newCtx, cancelFunc = context.WithTimeout(context.Background(), duration)
		}
	}
	return func() {
		thresholdFlags := ""
		code := http.StatusOK
		if cancelFunc != nil {
			cancelFunc()
		}
		if statusCode != nil {
			code = statusCode()
		}
		if code == runtime.StatusDeadlineExceeded {
			thresholdFlags = upstreamTimeoutFlag
		}
		access.Log(access.EgressTraffic, start, time.Since(start), req, &http.Response{StatusCode: code, Status: ""}, routeName, "", Milliseconds(duration), thresholdFlags)
	}
}
