package log

import (
	"github.com/felixge/httpsnoop"
	"github.com/go-ai-agent/core/runtime"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func AddRequestId(req *http.Request) string {
	if req == nil {
		return ""
	}
	id := req.Header.Get(runtime.XRequestId)
	if len(id) == 0 {
		uid, _ := uuid.NewUUID()
		id = uid.String()
		req.Header.Set(runtime.XRequestId, runtime.GetOrCreateRequestId(req))
	}
	return id
}

// Configure as last handler in chain
//middleware2.ControllerHttpHostMetricsHandler(mux, ""), status

// HttpHostMetricsHandler - handler for Http request metrics
func HttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		AddRequestId(r)
		m := httpsnoop.CaptureMetrics(appHandler, w, r)
		// log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		AnyAccess(IngressTraffic, start, time.Since(start), r, &http.Response{StatusCode: m.Code, ContentLength: m.Written}, -1, "")
	})
	return wrappedH
}
