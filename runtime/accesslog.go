package runtime

import (
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

// AccessLogHandler - template access log handler interface
type AccessLogHandler interface {
	Handle(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string)
}

type StdAccess struct{}

func (StdAccess) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	fmt.Println(FmtLog(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}

type LogAccess struct{}

func (LogAccess) Write(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, proxy, proxyThreshold, statusFlags string) {
	log.Println(FmtLog(traffic, start, duration, req, resp, routeName, timeout, rateLimit, rateBurst, rateThreshold, proxy, proxyThreshold, statusFlags))
}
