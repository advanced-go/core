package host

import (
	"fmt"
	"github.com/advanced-go/core/access"
	"net/http"
	"time"
)

const (
	Authorization   = "Authorization"
	upstreamTimeout = "UT"
)

type ServeHTTPFunc func(w http.ResponseWriter, r *http.Request)

func NewIntermediary(c1 ServeHTTPFunc, c2 ServeHTTPFunc) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrap := newWrapper(w)
		if c1 != nil {
			c1(wrap, r)
		}
		if wrap.statusCode == http.StatusOK && c2 != nil {
			c2(w, r)
		}
	}
}

func NewControllerIntermediary(routeName string, c2 ServeHTTPFunc) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c2 == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: componet 2 is nil")
			return
		}
		start := time.Now().UTC()
		wrap := newWrapper(w)
		c2(wrap, r)
		access.Log(access.EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: wrap.statusCode, ContentLength: wrap.written}, routeName, "", 0, "")
	}
}

/*
	func NewControllerIntermediary(duration, routeName string, c2 ServeHTTPFunc) ServeHTTPFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if c2 == nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error: componet 2 is nil")
				return
			}
			if len(duration) == 0 {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "error: duration is empty")
				return
			}
			threshold, err := controller.ParseDuration(duration)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", err)
				return
			}
			var writer *exchange.ResponseWriter
			if w2, ok := any(w).(*exchange.ResponseWriter); ok {
				writer = w2
			}
			thresholdFlags := ""
			start := time.Now().UTC()
			wrap := newWrapper(w)
			ch := make(chan struct{}, 1)
			ctx, cancelFunc := context.WithTimeout(context.Background(), threshold)
			defer cancelFunc()
			statusCode := http.StatusGatewayTimeout

			go func() {
				c2(wrap, r)
				ch <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				if writer != nil {
					writer.SetStatusCode(statusCode)
				}
				wrap.written = 0
				thresholdFlags = upstreamTimeout
			case <-ch:
				statusCode = wrap.statusCode
			}
			access.Log(access.EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: statusCode, ContentLength: wrap.written}, routeName, "", controller.Milliseconds(threshold), thresholdFlags)
		}
	}
*/
type wrapper struct {
	writer     http.ResponseWriter
	statusCode int
	written    int64
}

func newWrapper(writer http.ResponseWriter) *wrapper {
	w := new(wrapper)
	w.writer = writer
	w.statusCode = http.StatusOK
	return w
}

func (w *wrapper) Header() http.Header {
	return w.writer.Header()
}

func (w *wrapper) Write(p []byte) (int, error) {
	w.written += int64(len(p))
	return w.writer.Write(p)
}

func (w *wrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
}
