package messaging

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const (
	maxWait = time.Second * 3
)

const (
	pingLocation = PkgPath + ":Ping"
)

// Ping - templated function to "ping" a resource
func Ping(ctx context.Context, uri string) *Status {
	return ping(HostExchange, ctx, uri)
}

func ping(ex *Exchange, ctx context.Context, uri string) *Status {
	if uri == "" {
		return NewStatusError(errors.New("error: Ping() uri is empty"), pingLocation)
	}
	cache := NewMessageCache()
	msg := Message{To: uri, From: PkgPath, Event: PingEvent, ReplyTo: NewMessageCacheHandler(cache)}
	err := ex.SendCtrl(msg)
	if err != nil {
		return NewStatusError(err, pingLocation)
	}
	duration := maxWait
	for wait := time.Duration(float64(duration) * 0.20); duration >= 0; duration -= wait {
		time.Sleep(wait)
		result, ok := cache.Get(uri)
		if !ok {
			continue
		}
		if result.Status.Error != nil {
			return NewStatusError(errors.New(fmt.Sprintf("ping response status not available: [%v]", uri)), pingLocation)
		}
		return StatusOK()
	}
	return NewStatusError(errors.New(fmt.Sprintf("ping response time out: [%v]", uri)), pingLocation)
}
