package host

const (
	route = "host"
)

// Configure as last handler in chain
//middleware2.ControllerHttpHostMetricsHandler(mux, ""), status

// HttpHostMetricsHandler - handler for Http request metrics and access logging
/*
func HttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		m := httpsnoop.CaptureMetrics(appHandler, w, r)
		// log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		access.Log(access.IngressTraffic, start, time.Since(start), r, &http.Response{StatusCode: m.Code, ContentLength: m.Written}, route, "", -1, "")
	})
	return wrappedH
}


*/
