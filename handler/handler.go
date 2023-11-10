package handler

import (
	"github.com/felixge/httpsnoop"
	"github.com/go-ai-agent/core/httpx"
	"github.com/go-ai-agent/core/log"
	"net/http"
	"time"
)

// Configure as last handler in chain
//middleware2.ControllerHttpHostMetricsHandler(mux, ""), status

// HttpHostMetricsHandler - handler for Http request metrics
func HttpHostMetricsHandler[T log.AccessHandler](appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		r = httpx.UpdateHeadersAndContext(r)
		m := httpsnoop.CaptureMetrics(appHandler, w, r)
		var t T
		// log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		t(log.IngressTraffic, start, time.Since(start), r, &http.Response{StatusCode: m.Code, ContentLength: m.Written}, -1, "")
	})
	return wrappedH
}
