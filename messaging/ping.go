package messaging

import (
	"context"
	"errors"
	"fmt"
	"net/http"
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
	to, status := createTo(uri)
	if !status.OK() {
		return status
	}
	var response *Message

	duration := time.Second * 2
	result := make(chan *Status)
	reply := make(chan *Message, 16)
	msg := NewMessage(to, PkgPath, PingEvent)
	msg.ReplyTo = NewReceiverReplyTo(reply)
	err := ex.SendCtrl(msg)
	if err != nil {
		return NewStatusError(err, pingLocation)
	}
	go Receiver(duration, reply, result, func(msg *Message) bool {
		response = msg
		return true
	})
	status = <-result
	status.Location = pingLocation
	if response != nil {
		status.Code = response.Status.Code
		status.Error = response.Status.Error
	}
	// Not closing reply on a timeout as there is still an agent with a pending send on the reply channel.
	if status.Code != http.StatusGatewayTimeout {
		close(reply)
	}
	close(result)
	//close(reply)
	return status
}

func createTo(uri any) (string, *Status) {
	if uri == nil {
		return "", NewStatusError(errors.New("error: Ping() uri is nil"), pingLocation)
	}
	path := ""
	if u, ok := uri.(*url.URL); ok {
		path = u.Path
	} else {
		if u2, ok1 := uri.(string); ok1 {
			path = u2
		} else {
			return "", NewStatusError(errors.New(fmt.Sprintf("error: Ping() uri is invalid type: %v", reflect.TypeOf(uri).String())), pingLocation)
		}
	}
	nid, _, ok := UprootUrn(path)
	if !ok {
		return "", NewStatusError(errors.New(fmt.Sprintf("error: Ping() uri is not a valid URN %v", path)), pingLocation)
	}
	return nid, StatusOK()
}
