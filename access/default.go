package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func defaultFormatter(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) string {
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
		FmtJsonString(o.Region),
		FmtJsonString(o.Zone),
		FmtJsonString(o.SubZone),
		FmtJsonString(o.App),
		FmtJsonString(o.InstanceId),

		traffic,
		FmtTimestamp(start),
		strconv.Itoa(d),

		FmtJsonString(req.Header.Get(runtime.XRequestId)),
		FmtJsonString(req.Header.Get(runtime.XRelatesTo)),
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

func FmtJsonString(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "\"" + value + "\""
}

var defaultLogger = func(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	s := formatter(o, traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
	log.Default().Printf("%v\n", s)
}
