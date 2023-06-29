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

	Start          = "start"
	DurationMs     = "duration-ms"
	Traffic        = "traffic"
	Route          = "route"
	RequestId      = "request-id"
	Protocol       = "protocol"
	Method         = "method"
	Uri            = "uri"
	Host           = "host"
	Path           = "path"
	StatusCode     = "status-code"
	TimeoutMs      = "timeout-ms"
	RateLimit      = "rate-limit"
	RateBurst      = "rate-burst"
	RateThreshold  = "rate-threshold"
	Proxy          = "proxy"
	ProxyThreshold = "proxy-threshold"
	StatusFlags    = "status-flags"
)

// AccessLogHandler - template access log handler interface
type AccessLogHandler interface {
	Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string)
}

type StdioAccessLog struct{}

func (StdioAccessLog) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	fmt.Println(FormatLogJson(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}

type LogAccessLog struct{}

func (LogAccessLog) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	log.Println(FormatLogJson(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}

func FormatLogJson(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	sb := strings.Builder{}

	writeMarkup(&sb, Start, FmtTimestamp(start), true)
	writeMarkup(&sb, DurationMs, strconv.Itoa(d), false)
	writeMarkup(&sb, Traffic, traffic, true)
	writeMarkup(&sb, Route, routeName, true)
	writeMarkup(&sb, RequestId, req.Header.Get(RequestIdHeaderName), true)
	writeMarkup(&sb, Protocol, req.Proto, true)
	writeMarkup(&sb, Method, req.Method, true)
	writeMarkup(&sb, Uri, req.URL.String(), true)
	writeMarkup(&sb, Host, req.URL.Host, true)
	writeMarkup(&sb, Path, req.URL.Path, true)
	writeMarkup(&sb, StatusCode, strconv.Itoa(resp.StatusCode), false)

	writeMarkup(&sb, TimeoutMs, strconv.Itoa(timeout), false)

	writeMarkup(&sb, RateLimit, fmt.Sprintf("%v", rateLimit), false)
	writeMarkup(&sb, RateBurst, strconv.Itoa(rateBurst), false)
	writeMarkup(&sb, RateThreshold, rateThreshold, true)

	writeMarkup(&sb, Proxy, proxy, false)
	writeMarkup(&sb, ProxyThreshold, proxyThreshold, true)

	writeMarkup(&sb, StatusFlags, statusFlags, true)

	sb.WriteString("}")
	return sb.String()
}

func FormatLogText(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) string {
	d := int(duration / time.Duration(1e6))
	s := fmt.Sprintf(Start+":%v, "+
		DurationMs+":%v, "+
		Traffic+":%v, "+
		Route+":%v, "+
		RequestId+":%v, "+
		Protocol+":%v, "+
		Method+":%v, "+
		Uri+":%v, "+
		Host+":%v, "+
		Path+":%v, "+
		StatusCode+":%v, "+
		TimeoutMs+":%v, "+
		RateLimit+":%v, "+
		RateBurst+":%v, "+
		RateThreshold+":%v, "+
		Proxy+":%v, "+
		ProxyThreshold+":%v, "+
		StatusFlags+":%v",
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
