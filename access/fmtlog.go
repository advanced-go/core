package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	strings2 "github.com/advanced-go/core/strings"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func fmtLog(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, threshold int, thresholdFlags string) string {
	if req == nil {
		req, _ = http.NewRequest("", "https://somehost.com/search?q=test", nil)
	}
	if resp == nil {
		resp = &http.Response{StatusCode: http.StatusOK}
	}
	host := req.Host
	if len(host) == 0 {
		host = req.URL.Host
	}
	url := req.URL.String()
	if len(host) == 0 {
		//url = "urn:" + url
	} else {
		if len(req.URL.Scheme) == 0 {
			url = "http://" + host + req.URL.Path
		}
	}
	path := req.URL.Path
	i := strings.Index(path, ":")
	if i >= 0 {
		path = path[i+1:]
	}
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf("{ \"traffic\":\"%v\", "+
		"\"start\":%v, "+
		"\"duration\":%v, "+
		"\"route\":%v, "+
		"\"relates-to\":%v, "+
		"\"request-id\":%v, "+
		"\"protocol\":%v, "+
		"\"method\":%v, "+
		"\"uri\":%v, "+
		"\"host\":%v, "+
		"\"path\":%v, "+
		"\"status-code\":%v, "+
		"\"threshold\":%v, "+
		"\"threshold-flags\":%v }",
		traffic,
		strings2.FmtTimestamp(start),
		strconv.Itoa(d),
		fmtstr(routeName),
		fmtstr(req.Header.Get("RelatesTo")),
		fmtstr(req.Header.Get(runtime.XRequestId)),
		fmtstr(req.Proto),
		fmtstr(req.Method),
		fmtstr(url),  // fmtstr(req.URL.String()),
		fmtstr(host), // fmtstr(req.URL.Host),
		fmtstr(path),

		resp.StatusCode,

		threshold,
		fmtstr(thresholdFlags),
	)

	return s
}

func fmtstr(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "\"" + value + "\""
}
