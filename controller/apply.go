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
)

func Apply(ctx context.Context, newCtx *context.Context, req *http.Request, resp **http.Response, routeName string, duration time.Duration, statusCode access.StatusCodeFunc) func() {
	var cancelFunc context.CancelFunc

	if ctx == nil {
		ctx = context.Background()
	}
	*newCtx = ctx
	start := time.Now()
	// if a timeout and there is no deadline in the current ctx, then create a new context, otherwise update the duration with time
	// until the context deadline.
	if duration > 0 {
		if ct, ok := ctx.Deadline(); ok {
			duration = time.Until(ct)
		} else {
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
		if resp == nil || *resp == nil {
			r := new(http.Response)
			r.StatusCode = code
			r.Status = runtime.HttpStatus(code)
			if resp == nil {
				resp = &r
			} else {
				*resp = r
			}
		}
		access.Log(access.EgressTraffic, start, time.Since(start), req, *resp, routeName, "", Milliseconds(duration), thresholdFlags)
	}
}
