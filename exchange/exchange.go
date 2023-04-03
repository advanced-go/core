package exchange

import (
	"crypto/tls"
	"errors"
	"github.com/go-sre/core/runtime"
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
	if proxies, ok := runtime.IsProxyable(req.Context()); ok {
		do := findProxy(proxies)
		if do != nil {
			return do(req)
		}
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

func findProxy(proxies []any) func(*http.Request) (*http.Response, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*http.Request) (*http.Response, error)); ok {
			return fn
		}
	}
	return nil
}
