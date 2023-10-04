package exchange

import (
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

var (
	doLocation = PkgUrl + "/do"
	hdr        Default
)

func Do[E runtime.ErrorHandler](req *http.Request) (resp *http.Response, status *runtime.Status) {
	var e E

	if req == nil {
		return nil, e.Handle(nil, doLocation, errors.New("invalid argument : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	var err error
	resp, err = hdr.Do(req)
	if err != nil {
		return nil, e.Handle(req.Context(), doLocation, err).SetCode(http.StatusInternalServerError)
	}
	if resp == nil {
		return nil, e.Handle(req.Context(), doLocation, errors.New("invalid argument : response is nil")).SetCode(http.StatusInternalServerError)
	}
	return resp, runtime.NewHttpStatusCode(resp.StatusCode)
}

func DoT[E runtime.ErrorHandler, T any](req *http.Request) (resp *http.Response, t T, status *runtime.Status) {
	resp, status = Do[E](req)
	if !status.OK() {
		return nil, t, status
	}
	t, status = Deserialize[E, T](req.Context(), resp.Body)
	return
}
