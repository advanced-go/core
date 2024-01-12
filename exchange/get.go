package exchange

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	getLocation = PkgPath + ":Get"
)

func Get(uri string, h http.Header) (resp *http.Response, status runtime.Status) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, runtime.NewStatusError(http.StatusBadRequest, getLocation, err)
	}
	req.Header = h
	// exchange.Do() will always return a non nil *http.Response
	resp, status = Do(req)
	if !status.OK() {
		status.AddLocation(getLocation)
	}
	return
}
