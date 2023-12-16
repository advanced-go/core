package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	strings2 "github.com/advanced-go/core/strings"
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
		"\"route-name\":%v, "+
		"\"route-to\":%v, "+
		"\"threshold\":%v, "+
		"\"threshold-flags\":%v }",
		fmtstr(o.Region),
		fmtstr(o.Zone),
		fmtstr(o.SubZone),
		fmtstr(o.App),
		fmtstr(o.InstanceId),

		traffic,
		strings2.FmtTimestamp(start),
		strconv.Itoa(d),

		fmtstr(req.Header.Get(runtime.XRequestId)),
		fmtstr(req.Header.Get(runtime.XRelatesTo)),
		fmtstr(req.Proto),
		fmtstr(req.Method),
		fmtstr(url),
		fmtstr(host),
		fmtstr(path),

		resp.StatusCode,

		fmtstr(routeName),
		fmtstr(routeTo),
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

var defaultLogger = func(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	s := formatter(o, traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
	log.Default().Printf("%v\n", s)
}
