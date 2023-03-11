package exchange

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

// HttpExchange - interface for Http request/response interaction
type HttpExchange interface {
	Do(req *http.Request) (*http.Response, error)
}

type Default struct{}

func (Default) Do(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("request is nil") //NewStatus(StatusInvalidArgument, doLocation, errors.New("request is nil"))
	}
	if e, ok := Cast(req.Context()); ok {
		return e.Do(req)
	}
	return Client.Do(req)
}

var Client = http.DefaultClient

func init() {
	t, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		// Used clone instead of assignment due to presence of sync.Mutex fields
		var transport = t.Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		transport.MaxIdleConns = 200
		transport.MaxIdleConnsPerHost = 100
		Client = &http.Client{Transport: transport, Timeout: time.Second * 5}
	} else {
		Client = &http.Client{Transport: http.DefaultTransport, Timeout: time.Second * 5}
	}
}

func Cast(ctx context.Context) (HttpExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(HttpExchange); ok {
		return e, true
	}
	return nil, false
}
