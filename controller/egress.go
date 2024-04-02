package controller

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func doEgress(duration time.Duration, do func(*http.Request) (*http.Response, *runtime.Status), req *http.Request) (resp *http.Response, status *runtime.Status) {
	ch := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	go func() {
		defer func() {
			// Check for when a timeout is reached, the channel is closed, and there is a pending send
			if r := recover(); r != nil {
				fmt.Printf("recovered in controller.doEgress() : %v\n", r)
			}
		}()
		resp, status = do(req)
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			resp = &http.Response{StatusCode: http.StatusGatewayTimeout}
			status = runtime.NewStatusError(runtime.StatusDeadlineExceeded, ctx.Err())
		} else {
			resp = &http.Response{StatusCode: http.StatusInternalServerError}
			status = runtime.NewStatusError(http.StatusInternalServerError, ctx.Err())
		}
	case <-ch:
	}
	close(ch)
	return
}
