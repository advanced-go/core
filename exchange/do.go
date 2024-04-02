package exchange

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/advanced-go/core/controller"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	internalError           = "Internal Error"
	fileScheme              = "file"
	contextDeadlineExceeded = "context deadline exceeded"
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

func DeadlineExceededError(t any) bool {
	if t == nil {
		return false
	}
	if r, ok := t.(*http.Request); ok {
		return r.Context() != nil && r.Context().Err() == context.DeadlineExceeded
	}
	if e, ok := t.(error); ok {
		return e == context.DeadlineExceeded
	}
	return false
}

// Do - process a request, checking for overrides of file://, and a registered controller.
func Do(req *http.Request) (resp *http.Response, status *runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	if req.URL.Scheme == fileScheme {
		resp1, status1 := readResponse(req.URL)
		if !status1.OK() {
			return resp1, status1.AddLocation()
		}
		return resp1, runtime.NewStatus(resp1.StatusCode)
	}
	ctrl, status1 := controller.Lookup(req.URL.Path)
	if status1.OK() {
		return ctrl.Do(DoHttp, req)
	}
	return DoHttp(req)
}

// DoHttp - process an HTTP call
func DoHttp(req *http.Request) (resp *http.Response, status *runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	var err error

	resp, err = client.Do(req)
	if err != nil {
		// catch connectivity error, even with a valid URL
		if resp == nil {
			resp = serverErrorResponse()
		}
		// check for an error of deadline exceeded
		if req.Context() != nil && req.Context().Err() == context.DeadlineExceeded {
			resp.StatusCode = http.StatusGatewayTimeout
		}
		return resp, runtime.NewStatusError(resp.StatusCode, err)
	}
	return resp, runtime.NewStatus(resp.StatusCode)
}

func serverErrorResponse() *http.Response {
	resp := new(http.Response)
	resp.StatusCode = http.StatusInternalServerError
	resp.Status = internalError
	return resp
}
