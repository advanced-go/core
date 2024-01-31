package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func defaultFormatter(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) string {
	if req == nil {
		req, _ = http.NewRequest("", "https://somehost.com/search?q=test", nil)
	}
	if resp == nil {
		resp = &http.Response{StatusCode: http.StatusOK}
	}
	url, host, path := CreateUrlHostPath(req)
	s := fmt.Sprintf("{"+
		"\"region\":%v, "+
		"\"zone\":%v, "+
		"\"sub-zone\":%v, "+
		"\"app\":%v, "+
		"\"instance-id\":%v, "+
		" \"traffic\":\"%v\", "+
		"\"start\":%v, "+
		"\"duration\":%v, "+
		"\"request-id\":%v, "+
		"\"relates-to\":%v, "+
		"\"protocol\":%v, "+
		"\"method\":%v, "+
		"\"uri\":%v, "+
		"\"host\":%v, "+
		"\"path\":%v, "+
		"\"status-code\":%v, "+
		"\"status\":%v, "+
		"\"route\":%v, "+
		"\"route-to\":%v, "+
		"\"threshold\":%v, "+
		"\"threshold-flags\":%v }",
		FmtJsonString(""),
		FmtJsonString(""),
		FmtJsonString(""),
		FmtJsonString(""),
		FmtJsonString(""),

		traffic,
		FmtTimestamp(start),
		strconv.Itoa(Milliseconds(duration)),

		FmtJsonString(req.Header.Get(xRequestId)),
		FmtJsonString(req.Header.Get(xRelatesTo)),
		FmtJsonString(req.Proto),
		FmtJsonString(req.Method),
		FmtJsonString(url),
		FmtJsonString(host),
		FmtJsonString(path),

		resp.StatusCode,
		FmtJsonString(resp.Status),

		FmtJsonString(routeName),
		FmtJsonString(routeTo),
		threshold,
		FmtJsonString(thresholdFlags),
	)

	return s
}

var defaultLogger = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	s := defaultFormatter(traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
	fmt.Printf("%v\n", s)
}
