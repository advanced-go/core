package exchange

import (
	"crypto/tls"
	"errors"
	"github.com/go-ai-agent/core/httpx"
	"github.com/go-ai-agent/core/runtime"
	"google.golang.org/grpc/codes"
	"net/http"
	"time"
)

// HttpExchange - interface for Http request/response interaction
//type HttpExchange interface {
//	Do(req *http.Request) (*http.Response, error)
//}

var (
	doLocation = PkgUri + "/Do"
	Client     = http.DefaultClient
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

func Do(req *http.Request) (resp *http.Response, status *runtime.Status) {
	if req == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, doLocation, errors.New("invalid argument : request is nil")) //.SetCode(runtime.StatusInvalidArgument)
	}
	var err error
	var doProxy Exchange

	if runtime.IsDebugEnvironment() {
		if req.URL.Scheme == "file" {
			resp, err = httpx.ReadResponse(req.URL)
			if err != nil {
				return resp, runtime.NewStatusError(http.StatusInternalServerError, doLocation, err)
			}
			return resp, runtime.NewStatusOK()
		}
		if proxies, ok := runtime.IsProxyable(req.Context()); ok {
			do := findProxy(proxies)
			if do != nil {
				doProxy = do
			}
		}
	}
	if doProxy == nil {
		doProxy = ProxyLookup(req)
	}
	if doProxy != nil {
		resp, err = doProxy(req)
	} else {
		resp, err = Client.Do(req)
	}
	if err != nil {
		if resp == nil {
			resp = &http.Response{StatusCode: http.StatusInternalServerError, Status: "internal server error"}
		}
		return resp, runtime.NewStatusError(codes.Code(resp.StatusCode), doLocation, err)
	}
	return resp, runtime.NewHttpStatus(resp.StatusCode)
}

func DoT[T any](req *http.Request) (resp *http.Response, t T, status *runtime.Status) {
	resp, status = Do(req)
	if !status.OK() {
		return nil, t, status
	}
	t, status = Deserialize[T](resp.Body)
	//var e E
	//if !status.OK() {
	//	e.HandleStatus(status, req) //.SetRequestId(runtime.RequestId(req)))
	//}
	return
}

func findProxy(proxies []any) func(*http.Request) (*http.Response, error) {
	for _, p := range proxies {
		if fn, ok := p.(func(*http.Request) (*http.Response, error)); ok {
			return fn
		}
	}
	return nil
}
