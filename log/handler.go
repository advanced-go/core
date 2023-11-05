package log

//middleware2.ControllerHttpHostMetricsHandler(mux, ""), status
// middleware2.ControllerWrapTransport(exchange.Client)

// ControllerHttpHostMetricsHandler - handler that applies host and ingress controllers
/*
func ControllerHttpHostMetricsHandler(appHandler http.Handler, msg string) http.Handler {
	wrappedH := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now().UTC()
		ctrl := controller.IngressLookupHost()
		var m httpsnoop.Metrics

		if ctrl != nil {
			if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
				controller.LogHttpIngress(ctrl, start, time.Since(start), r, rlc.StatusCode(), 0, controller.RateLimitFlag)
				return
			}
		}
		ctrl, _ = controller.IngressLookup(r)
		if ctrl == nil {
			m = httpsnoop.CaptureMetrics(appHandler, w, r)
		} else {
			if toc := ctrl.Timeout(); toc.IsEnabled() && toc.Duration() > 0 {
				m = httpsnoop.CaptureMetrics(http.TimeoutHandler(appHandler, toc.Duration(), msg), w, r)
			} else {
				m = httpsnoop.CaptureMetrics(appHandler, w, r)
			}
		}
		// log.Printf("%s %s (code=%d dt=%s written=%d)", r.Method, r.URL, m.Code, m.Duration, m.Written)
		// TO DO: determine how to set status flag value when a timeout occurs
		controller.LogHttpIngress(ctrl, start, time.Since(start), r, m.Code, m.Written, "")
	})
	return wrappedH
}


*/
