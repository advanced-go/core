package http2

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/advanced-go/core/http2/io"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	doLocation     = PkgPath + ":Do"
	doReadLocation = PkgPath + ":Do/readFromFile"
)

var (
	Client = http.DefaultClient
)

func init() {
	fmt.Println("do.go -> init()")
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

	if req.URL.Scheme == FileScheme {
		resp1, err1 := io.ReadResponse(req.URL)
		if err1 != nil {
			if resp1 == nil {
				resp1 = new(http.Response)
				resp1.StatusCode = http.StatusInternalServerError
				resp1.Status = "internal server error"
			}
			return resp1, runtime.NewStatusError(http.StatusInternalServerError, doReadLocation, err1)
		}
		return resp1, runtime.NewStatus(resp1.StatusCode)
	}
	resp, err = Client.Do(req)
	if err != nil {
		// Happens as a result of a connectivity error, even with a valid URL
		if resp == nil {
			resp = new(http.Response)
			resp.StatusCode = http.StatusInternalServerError
			resp.Status = "internal server error"
		}
		return resp, runtime.NewStatusError(resp.StatusCode, doLocation, err)
	}
	return resp, runtime.NewStatus(resp.StatusCode)
}

// DoT - do a Http exchange with deserialization
func DoT[T any](req *http.Request) (resp *http.Response, t T, status runtime.Status) {
	resp, status = Do(req)
	if !status.OK() {
		return nil, t, status
	}
	t, status = Deserialize[T](resp.Body)
	return
}

/*
func readFromFile(req *http.Request) (*http.Response, runtime.Status) {
	var uri *url.URL
	var err error


		if req.URL.Scheme == FileScheme {
			uri = req.URL
		} else {
			location := req.Header.Get(ContentLocation)
			if len(location) == 0 {
				return nil, runtime.StatusOK()
			}
			uri, err = url.Parse(location)
			if err != nil {
				return nil, runtime.NewStatusError(http.StatusInternalServerError, doReadLocation, err)
			}
		}
		if uri == nil {
			return nil, runtime.StatusOK()
		}
	resp, err1 := io.ReadResponse(req.URL)
	if err1 != nil {
		return resp, runtime.NewStatusError(http.StatusInternalServerError, doReadLocation, err1)
	}
	return resp, runtime.NewStatus(resp.StatusCode)
}

*/
