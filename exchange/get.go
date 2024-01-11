package exchange

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	getLocation = PkgPath + ":Get"
)

func Get[E runtime.ErrorHandler](uri string, h http.Header, callerLocation string) (*http.Response, runtime.Status) {
	var e E

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, e.Handle(runtime.NewStatusError(http.StatusInternalServerError, getLocation, err), runtime.RequestId(h), callerLocation)
	}
	req.Header = h
	// exchange.Do() will always return a non nil *http.Response
	resp, status := Do(req)
	if !status.OK() {
		return nil, e.Handle(status, runtime.RequestId(h), callerLocation)
	}
	return resp, runtime.StatusOK()
}
