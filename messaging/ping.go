package messaging

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"time"
)

const (
	maxWait = time.Second * 3
)

const (
	pingLocation = PkgPath + ":Ping"
)

// Ping - function to "ping" a resource
func Ping(ctx context.Context, uri any) *Status {
	return ping(ctx, HostExchange, uri)
}

func ping(ctx context.Context, ex *Exchange, uri any) *Status {
	if uri == nil {
		return NewStatusError(errors.New("error: Ping() uri is nil"), pingLocation)
	}
	path := ""
	if u, ok := uri.(*url.URL); ok {
		path = u.Path
	} else {
		if u2, ok1 := uri.(string); ok1 {
			path = u2
		} else {
			return NewStatusError(errors.New(fmt.Sprintf("error: Ping() uri is invalid type: %v", reflect.TypeOf(uri).String())), pingLocation)
		}
	}
	nid, _, ok := UprootUrn(path)
	if !ok {
		return NewStatusError(errors.New(fmt.Sprintf("error: Ping() uri is not a valid URN %v", path)), pingLocation)
	}
	cache := NewMessageCache()
	msg := Message{To: nid, From: PkgPath, Event: PingEvent, ReplyTo: NewMessageCacheHandler(cache)}
	err := ex.SendCtrl(msg)
	if err != nil {
		return NewStatusError(err, pingLocation)
	}
	duration := maxWait
	for wait := time.Duration(float64(duration) * 0.20); duration >= 0; duration -= wait {
		time.Sleep(wait)
		result, ok2 := cache.Get(nid)
		if !ok2 {
			continue
		}
		if result.Status.Error != nil {
			return NewStatusError(errors.New(fmt.Sprintf("ping response status not available: [%v]", uri)), pingLocation)
		}
		return StatusOK()
	}
	return NewStatusError(errors.New(fmt.Sprintf("ping response time out: [%v]", uri)), pingLocation)
}
