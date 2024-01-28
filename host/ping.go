package host

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/messaging"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	maxWait = time.Second * 3
)

const (
	pingLocation = PkgPath + ":Ping"
)

// Ping - templated function to "ping" a resource
func Ping[E runtime.ErrorHandler](ctx context.Context, uri string) (status runtime.Status) {
	return ping[E](messaging.HostExchange, ctx, uri)
}

func ping[E runtime.ErrorHandler](ex messaging.Exchange, ctx context.Context, uri string) (status runtime.Status) {
	var e E

	if uri == "" {
		return e.Handle(runtime.NewStatusError(runtime.StatusInvalidArgument, pingLocation, errors.New("invalid argument: ping uri is empty")),
			runtime.RequestId(ctx), "")
	}
	cache := messaging.NewMessageCache()
	msg := messaging.Message{To: uri, From: PkgPath, Event: messaging.PingEvent, Status: nil, ReplyTo: messaging.NewMessageCacheHandler[E](cache)}
	status = ex.SendCtrl(msg)
	if !status.OK() {
		return e.Handle(status, runtime.RequestId(ctx), pingLocation)
	}
	duration := maxWait
	for wait := time.Duration(float64(duration) * 0.20); duration >= 0; duration -= wait {
		time.Sleep(wait)
		result, err1 := cache.Get(uri)
		if err1 != nil {
			//fmt.Printf("wait: [%v] error: [%v]\n", wait, err1)
			continue
		}
		if result.Status == nil {
			return e.Handle(runtime.NewStatusError(http.StatusInternalServerError, pingLocation, errors.New(fmt.Sprintf("ping response status not available: [%v]", uri))), runtime.RequestId(ctx), "")
		}
		return result.Status
	}
	return e.Handle(runtime.NewStatusError(runtime.StatusDeadlineExceeded, pingLocation, errors.New(fmt.Sprintf("ping response time out: [%v]", uri))), runtime.RequestId(ctx), "")
}
