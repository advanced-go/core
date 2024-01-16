package exchange

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	getLocation = PkgPath + ":Get"
)

// Get - process an HTTP Get request
func Get(uri string, h http.Header) (resp *http.Response, status runtime.Status) {
	if len(uri) == 0 {
		return serverErrorResponse(), runtime.NewStatusError(http.StatusBadRequest, getLocation, errors.New("error: URI is empty"))
	}
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return serverErrorResponse(), runtime.NewStatusError(http.StatusBadRequest, getLocation, err)
	}
	req.Header = h
	// exchange.Do() will always return a non nil *http.Response
	resp, status = Do(req)
	if !status.OK() {
		status.AddLocation(getLocation)
	}
	return
}
