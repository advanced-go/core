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

	start := time.Now()
	req, _ := http.NewRequest(method, uri, nil)
	if h != nil {
		req.Header = h
	}
	// TO DO : determine if the current context already contains a CancelCtx
	if duration > 0 {
		if ctx != nil {
		} else {
			*newCtx, cancelFunc = context.WithTimeout(context.Background(), duration)
		}
	} else {
		if ctx == nil {
			*newCtx = context.Background()
		} else {
			*newCtx = ctx
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
