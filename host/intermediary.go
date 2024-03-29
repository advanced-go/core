package host

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/controller"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	Authorization = "Authorization"
	Timeout       = "TO"
	XRequestId    = "X-Request-Id"
)

type ServeHTTPFunc func(w http.ResponseWriter, r *http.Request)

func NewConditionalIntermediary(c1 ServeHTTPFunc, c2 ServeHTTPFunc, ok func(int) bool) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c2 == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: component 2 is nil")
			return
		}
		w2 := newWrapper(w)
		if c1 != nil {
			c1(w2, r)
		}
		if (ok == nil && w2.statusCode == http.StatusOK) || (ok != nil && ok(w2.statusCode)) {
			c2(w, r)
		}
	}
}

func NewControllerIntermediary2(routeName string, c2 ServeHTTPFunc) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c2 == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: component 2 is nil")
			return
		}
		start := time.Now().UTC()
		wrap := newWrapper(w)
		c2(wrap, r)
		access.Log(access.EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: wrap.statusCode, ContentLength: wrap.written}, routeName, "", 0, "")
	}
}

func NewControllerIntermediary(ctrl *controller.Control2, c2 ServeHTTPFunc) ServeHTTPFunc {
	return newControllerIntermediary(ctrl, c2, access.InternalTraffic)
}

func newControllerIntermediary(ctrl *controller.Control2, c2 ServeHTTPFunc, traffic string) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c2 == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: component 2 is nil")
			return
		}
		if traffic == access.IngressTraffic {
			if r.Header.Get(XRequestId) == "" {
				uid, _ := uuid.NewUUID()
				r.Header.Add(XRequestId, uid.String())
			}
		}
		routeName := ""
		flags := ""
		var start time.Time
		var duration time.Duration
		if ct, ok := r.Context().Deadline(); ok {
			duration = time.Until(ct) * -1
		}
		w2 := newWrapper(w)
		if ctrl != nil && ctrl.Timeout.Duration > 0 && duration == 0 {
			routeName = ctrl.RouteName
			duration = ctrl.Timeout.Duration
			ctx, cancel := context.WithTimeout(r.Context(), ctrl.Timeout.Duration)
			defer cancel()
			r2 := r.Clone(ctx)
			start = time.Now().UTC()
			c2(w2, r2)
		} else {
			start = time.Now().UTC()
			c2(w2, r)
		}
		if w2.statusCode == http.StatusGatewayTimeout {
			flags = Timeout
		}
		if traffic == "" {
			traffic = access.InternalTraffic
		}
		access.Log(traffic, start, time.Since(start), r, &http.Response{StatusCode: w2.statusCode, ContentLength: w2.written}, routeName, "", Milliseconds(duration), flags)
	}
}

func NewAccessLogIntermediary(routeName string, c2 ServeHTTPFunc) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c2 == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error: component 2 is nil")
			return
		}
		w2 := newWrapper(w)
		start := time.Now().UTC()
		c2(w2, r)
		access.Log(access.InternalTraffic, start, time.Since(start), r, &http.Response{StatusCode: w2.statusCode, ContentLength: w2.written}, routeName, "", 0, "")
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

// Milliseconds - convert time.Duration to milliseconds
func Milliseconds(duration time.Duration) int {
	return int(duration / time.Duration(1e6))
}
