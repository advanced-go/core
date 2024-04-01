package host

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	RouteName = "host"
)

func NewHostTimeoutIntermediary(d time.Duration, c2 HttpHandlerFunc) HttpHandlerFunc {
	return newIngressTimeoutIntermediary(RouteName, d, c2, access.IngressTraffic)
}

func NewIngressTimeoutIntermediary(routeName string, d time.Duration, c2 HttpHandlerFunc) HttpHandlerFunc {
	return newIngressTimeoutIntermediary(routeName, d, c2, access.InternalTraffic)
}

func newIngressTimeoutIntermediary(routeName string, d time.Duration, c2 HttpHandlerFunc, traffic string) HttpHandlerFunc {
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
		w2 := newWrapper(w)
		apply(w2, r, routeName, d, c2, traffic, "")
	}
}

func apply(w *wrapper, r *http.Request, routeName string, duration time.Duration, handler HttpHandlerFunc, traffic, routeTo string) {
	if handler == nil {
		return
	}
	existing := false
	flags := ""
	var start time.Time
	if ct, ok := r.Context().Deadline(); ok {
		existing = true
		duration = time.Until(ct) * -1
	}
	if !existing && duration > 0 {
		ctx, cancel := context.WithTimeout(r.Context(), duration)
		defer cancel()
		r2 := r.Clone(ctx)
		start = time.Now().UTC()
		handler(w, r2)
	} else {
		start = time.Now().UTC()
		handler(w, r)
	}
	if w.statusCode == http.StatusGatewayTimeout {
		flags = TimeoutFlag
	}
	if traffic == "" {
		traffic = access.InternalTraffic
	}
	access.Log(traffic, start, time.Since(start), r, &http.Response{StatusCode: w.statusCode, ContentLength: w.written}, routeName, routeTo, Milliseconds(duration), flags)
}

/*
	func NewTimeoutIntermediary(duration, routeName string, c2 HttpHandlerFunc) HttpHandlerFunc {
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

/*
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
		flags = TimeoutFlag
	}
	if traffic == "" {
		traffic = access.InternalTraffic
	}
	access.Log(traffic, start, time.Since(start), r, &http.Response{StatusCode: w2.statusCode, ContentLength: w2.written}, routeName, "", Milliseconds(duration), flags)
*/
