package access

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	ContentEncoding = "Content-Encoding"
)

func DefaultFormatter(o *Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) string {
	if o == nil {
		o = &origin
	}
	req = SafeRequest(req)
	resp = SafeResponse(resp)
	url, host, path := CreateUrlHostPath(req)
	s := fmt.Sprintf("{"+
		"\"region\":%v, "+
		"\"zone\":%v, "+
		"\"sub-zone\":%v, "+
		"\"app\":%v, "+
		"\"instance-id\":%v, "+
		"\"traffic\":\"%v\", "+
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
		"\"encoding\":%v, "+
		"\"bytes-written\":%v, "+
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
		strconv.Itoa(Milliseconds(duration)),

		FmtJsonString(req.Header.Get(XRequestId)),
		FmtJsonString(req.Header.Get(XRelatesTo)),
		FmtJsonString(req.Proto),
		FmtJsonString(req.Method),
		FmtJsonString(url),
		FmtJsonString(host),
		FmtJsonString(path),

		resp.StatusCode,
		//FmtJsonString(resp.Status),
		FmtJsonString(Encoding(resp)),
		fmt.Sprintf("%v", resp.ContentLength),

		FmtJsonString(routeName),
		FmtJsonString(routeTo),
		threshold,
		FmtJsonString(thresholdFlags),
	)

	return s
}

var defaultLogger = func(o *Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	s := formatter(o, traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
	log.Default().Printf("%v\n", s)
}
