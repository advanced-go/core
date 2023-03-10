package exchange

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

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

// Do - process a "client.Do" request with the http.DefaultClient
func Do(req *http.Request) (*http.Response, error) {
	return DoClient(req, Client)
}

// DoClient - process a "client.Do" request with the given client. Also, check the req.Context to determine
// if there is an Exchange interface that should be called instead.
func DoClient(req *http.Request, client *http.Client) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("request is nil") //NewStatus(StatusInvalidArgument, doLocation, errors.New("request is nil"))
	}
	if client == nil {
		return nil, errors.New("client is nil")
	}
	if e, ok := exchangeCast(req.Context()); ok {
		return e.Do(req)
	}
	return client.Do(req)
}

func exchangeCast(ctx context.Context) (HttpExchange, bool) {
	if ctx == nil {
		return nil, false
	}
	if e, ok := any(ctx).(HttpExchange); ok {
		return e, true
	}
	return nil, false
}
