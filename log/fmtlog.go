package log

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	strings2 "github.com/go-ai-agent/core/strings"
	"net/http"
	"strconv"
	"time"
)

//var LogFn = defaultLogFn //

var AccessLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string) {
	s := fmtLog(traffic, start, duration, req, resp, statusFlags)
	fmt.Printf("{%v}\n", s)
}

func fmtLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf("{ \"start\":%v, "+
		"\"duration\":%v, "+
		"\"traffic\":\"%v\", "+
		//"controller:%v, "+
		"\"request-id\":\"%v\", "+
		"\"protocol\":\"%v\", "+
		"\"method\":\"%v\", "+
		"\"url\":\"%v\", "+
		"\"host\":\"%v\", "+
		"\"path\":\"%v\", "+
		"\"status-code\":%v, "+
		//"timeout-ms:%v, "+
		//"rate-limit:%v, "+
		//"rate-burst:%v, "+
		//	"proxy:%v, "+
		"\"status-flags\":\"%v\" }",
		strings2.FmtTimestamp(start), //l.Value(StartTimeOperator),
		strconv.Itoa(d),              //l.Value(DurationOperator),
		traffic,                      //l.Value(TrafficOperator),
		//controllerName,

		req.Header.Get(runtime.XRequestId), //l.Value(RequestIdOperator),
		req.Proto,                          //l.Value(RequestProtocolOperator),
		req.Method,                         //l.Value(RequestMethodOperator),
		req.URL.String(),                   //l.Value(RequestUrlOperator),
		req.URL.Host,                       //l.Value(RequestHostOperator),
		req.URL.Path,                       //l.Value(RequestPathOperator),

		resp.StatusCode, //l.Value(ResponseStatusCodeOperator),

		//timeout, //Tl.Value(TimeoutDurationOperator),
		//rateLimit, //l.Value(RateLimitOperator),
		//rateBurst, //l.Value(RateBurstOperator),
		//proxy,
		statusFlags, //l.Value(StatusFlagsOperator),
	)

	return s
}
