package exchange

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	doLocation     = PkgPath + ":do"
	doReadResponse = PkgPath + ":readResponse"
	internalError  = "Internal Error"
	readStatusLoc  = PkgPath + ":readStatus"
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

func do(req *http.Request) (resp *http.Response, status runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, doLocation, errors.New("invalid argument : request is nil"))
	}
	var err error

	if uri.IsFileScheme(req.URL) {
		if uri.IsStatusURL(req.URL.String()) {
			return readStatus(req.URL)
		}
		resp1, status1 := readResponse(req.URL)
		if !status1.OK() {
			return resp1, status1.AddLocation(doLocation)
		}
		return resp1, runtime.NewStatus(resp1.StatusCode)
	}
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

func readStatus(u *url.URL) (*http.Response, runtime.Status) {
	if u != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusInvalidArgument, readStatusLoc, errors.New("invalid argument : URL is nil"))
	}
	buf, err1 := os.ReadFile(uri.FileName(u))
	if err1 != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusIOError, readStatusLoc, err1)
	}
	var status runtime.SerializedStatusState
	err := json.Unmarshal(buf, &status)
	if err != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusJsonDecodeError, readStatusLoc, err)
	}
	if len(status.Err) > 0 {
		return serverErrorResponse(), runtime.NewStatusError(status.Code, status.Location, errors.New(status.Err))
	}
	resp := new(http.Response)
	resp.StatusCode = status.Code
	return resp, runtime.NewStatus(status.Code).AddLocation(status.Location)
}
