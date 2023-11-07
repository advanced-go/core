package log

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	strings2 "github.com/go-ai-agent/core/strings"
	"net/http"
	"strconv"
	"time"
)

func fmtLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) string {
	if req == nil {
		req, _ = http.NewRequest("", "https://somehost.com/search?q=test", nil)
	}
	if resp == nil {
		resp = &http.Response{StatusCode: http.StatusOK}
	}
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf("{ \"start\":%v, "+
		"\"duration\":%v, "+
		"\"traffic\":\"%v\", "+
		//"controller:%v, "+
		"\"request-id\":%v, "+
		"\"protocol\":%v, "+
		"\"method\":%v, "+
		"\"url\":%v, "+
		"\"host\":%v, "+
		"\"path\":%v, "+
		"\"status-code\":%v, "+
		"\"threshold\":%v, "+
		"\"status-flags\":%v }",
		strings2.FmtTimestamp(start),
		strconv.Itoa(d),
		traffic,

		fmtstr(req.Header.Get(runtime.XRequestId)),
		fmtstr(req.Proto),
		fmtstr(req.Method),
		fmtstr(req.URL.String()),
		fmtstr(req.URL.Host),
		fmtstr(req.URL.Path),

		resp.StatusCode,

		threshold,
		fmtstr(statusFlags),
	)

	return s
}

func fmtstr(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "\"" + value + "\""
}
