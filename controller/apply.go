package controller

import (
	"context"
	"net/http"
	"time"
)

const (
	upstreamTimeoutFlag    = "UT"
	statusDeadlineExceeded = 4
	egressTraffic          = "egress"
	xRequestId             = "x-request-id"
	xRelatesTo             = "x-relates-to"
)

var (
	logger = defaultLogger
)

// Logger - log function
type Logger func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, routeTo string, threshold int, thresholdFlags string)

// SetLogger - override logging
func SetLogger(fn Logger) {
	if fn != nil {
		logger = fn
	}
}

type StatusCodeFunc func() int

func Apply(ctx context.Context, newCtx *context.Context, method, uri, routeName string, h http.Header, duration time.Duration, statusCode StatusCodeFunc) func() {
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
		threshold := Milliseconds(duration)
		if statusCode != nil {
			code = statusCode()
		}
		if code == statusDeadlineExceeded {
			thresholdFlags = upstreamTimeoutFlag
		} else {
			threshold = -1
		}
		logger(egressTraffic, start, time.Since(start), req, &http.Response{StatusCode: code, Status: ""}, routeName, "", threshold, thresholdFlags)
	}
}
