package exchange

import (
	"crypto/tls"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"time"
)

const (
	doLocation    = PkgPath + ":Do"
	internalError = "Internal Error"
)

var (
	client = http.DefaultClient
)

func init() {
	t, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		// Used clone instead of assignment due to presence of sync.Mutex fields
		var transport = t.Clone()
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		transport.MaxIdleConns = 200
		transport.MaxIdleConnsPerHost = 100
		client = &http.Client{Transport: transport, Timeout: time.Second * 5}
	} else {
		client = &http.Client{Transport: http.DefaultTransport, Timeout: time.Second * 5}
	}
}

// Do - process a request, checking for overrides of file://, and a registered endpoint.
func Do(req *http.Request) (resp *http.Response, status runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, doLocation, errors.New("invalid argument : request is nil"))
	}
	if uri.IsFileScheme(req.URL) {
		resp1, status1 := readResponse(req.URL)
		if !status1.OK() {
			return resp1, status1.AddLocation(doLocation)
		}
		return resp1, runtime.NewStatus(resp1.StatusCode)
	}
	handler, status1 := proxyLookup(req.URL)
	if status1.OK() {
		w := newResponseWriter()
		handler(w, req)
		resp = w.Result()
		return resp, runtime.NewStatus(resp.StatusCode)
	}
	return DoHttp(req)
}

// DoHttp - process and HTTP call
func DoHttp(req *http.Request) (resp *http.Response, status runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, doLocation, errors.New("invalid argument : request is nil"))
	}
	var err error

	//if uri.IsFileScheme(req.URL) {
	////	resp1, status1 := readResponse(req.URL)
	//	if !status1.OK() {
	//		return resp1, status1.AddLocation(doLocation)
	//	}
	//	return resp1, runtime.NewStatus(resp1.StatusCode)
	//}
	resp, err = client.Do(req)
	if err != nil {
		// catch connectivity error, even with a valid URL
		if resp == nil {
			resp = serverErrorResponse()
		}
		return resp, runtime.NewStatusError(resp.StatusCode, doLocation, err)
	}
	return resp, runtime.NewStatus(resp.StatusCode)
}

func serverErrorResponse() *http.Response {
	resp := new(http.Response)
	resp.StatusCode = http.StatusInternalServerError
	resp.Status = internalError
	return resp
}
