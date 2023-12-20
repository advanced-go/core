package http2test

import (
	"crypto/tls"
	"errors"
	"github.com/advanced-go/core/exchange"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"time"
)

const (
	doLocation     = ":Do"
	doReadResponse = ":readReponse"
	internalError  = "Internal Error"
)

var (
	Client = http.DefaultClient
)

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

// Do - do a Http exchange with a runtime.Status
func Do(req *http.Request) (resp *http.Response, status runtime.Status) {
	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, doLocation, errors.New("invalid argument : request is nil")) //.SetCode(runtime.StatusInvalidArgument)
	}
	var err error

	if uri.IsFileScheme(req.URL) {
		resp1, err1 := exchange.ReadResponse(req.URL)
		if err1 != nil {
			if resp1 == nil {
				resp1 = new(http.Response)
				resp1.StatusCode = http.StatusInternalServerError
				resp1.Status = internalError
			}
			return resp1, runtime.NewStatusError(http.StatusInternalServerError, doReadResponse, err1)
		}
		return resp1, runtime.NewStatus(resp1.StatusCode)
	}
	resp, err = Client.Do(req)
	if err != nil {
		// catch connectivity error, even with a valid URL
		if resp == nil {
			resp = new(http.Response)
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = "internal server error"
		}
		return resp, runtime.NewStatusError(resp.StatusCode, doLocation, err)
	}
	return resp, runtime.NewStatus(resp.StatusCode)
}
