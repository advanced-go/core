package runtime

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	markupNull          = "\"%v\":null"
	markupString        = "\"%v\":\"%v\""
	markupValue         = "\"%v\":%v"
	RequestIdHeaderName = "X-REQUEST-ID"
)

// AccessLogHandler - template access log handler interface
type AccessLogHandler interface {
	Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string)
}

type StdioWriter struct{}

func (StdioWriter) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	fmt.Println(FormatLogJson(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}

type LogWriter struct{}

func (LogWriter) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	log.Println(FormatLogJson(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}

func FormatLogJson(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	sb := strings.Builder{}

	writeMarkup(&sb, "start", FmtTimestamp(start), true)
	writeMarkup(&sb, "duration-ms", strconv.Itoa(d), false)
	writeMarkup(&sb, "traffic", traffic, true)
	writeMarkup(&sb, "route", routeName, true)
	writeMarkup(&sb, "request-id", req.Header.Get(RequestIdHeaderName), true)
	writeMarkup(&sb, "protocol", req.Proto, true)
	writeMarkup(&sb, "method", req.Method, true)
	writeMarkup(&sb, "uri", req.URL.String(), true)
	writeMarkup(&sb, "host", req.URL.Host, true)
	writeMarkup(&sb, "path", req.URL.Path, true)
	writeMarkup(&sb, "status-code", strconv.Itoa(resp.StatusCode), false)

	writeMarkup(&sb, "timeout-ms", strconv.Itoa(timeout), false)

	writeMarkup(&sb, "rate-limit", fmt.Sprintf("%v", rateLimit), false)
	writeMarkup(&sb, "rate-burst", strconv.Itoa(rateBurst), false)
	writeMarkup(&sb, "rate-threshold", rateThreshold, true)

	writeMarkup(&sb, "proxy", proxy, false)
	writeMarkup(&sb, "proxy-threshold", proxyThreshold, true)

	writeMarkup(&sb, "status-flags", statusFlags, true)

	sb.WriteString("}")
	return sb.String()
}

func FormatLogText(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf("start:%v, "+
		"duration-ms:%v, "+
		"traffic:%v, "+
		"route:%v, "+
		"request-id:%v, "+
		"protocol:%v, "+
		"method:%v, "+
		"url:%v, "+
		"host:%v, "+
		"path:%v, "+
		"status-code:%v, "+
		"timeout-ms:%v, "+
		"rate-limit:%v, "+
		"rate-burst:%v, "+
		"rate-threshold:%v, "+
		"proxy:%v, "+
		"proxy-threshold:%v, "+
		"status-flags:%v",
		FmtTimestamp(start),
		strconv.Itoa(d),
		traffic,
		routeName,

		req.Header.Get(RequestIdHeaderName),
		req.Proto,
		req.Method,
		req.URL.String(),
		req.URL.Host,
		req.URL.Path,

		resp.StatusCode,

		timeout,

		rateLimit,
		rateBurst,
		rateThreshold,

		proxy,
		proxyThreshold,
		statusFlags,
	)

	return s
}

func writeMarkup(sb *strings.Builder, name, value string, stringValue bool) {
	if sb.Len() == 0 {
		sb.WriteString("{")
	} else {
		sb.WriteString(",")
	}
	if len(value) == 0 {
		sb.WriteString(fmt.Sprintf(markupNull, name))
	} else {
		format := markupString
		if !stringValue {
			format = markupValue
		}
		sb.WriteString(fmt.Sprintf(format, name, value))
	}
}
