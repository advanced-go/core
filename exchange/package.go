package exchange

import (
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/url"
)

type pkg struct{}

const (
	PkgPath     = "github.com/advanced-go/core/exchange"
	doRouteName = "http-exchange"
	doLoc       = PkgPath + ":Do"
)

// Do - process a Http exchange with a runtime.Status
func Do(req *http.Request) (resp *http.Response, status runtime.Status) {
	r := access.NewRequest(req.Header, req.Method, doLoc)
	defer access.LogDeferred(access.InternalTraffic, r, doRouteName, "", -1, "", &status)()
	return do(req)
}

// ReadResponse - read a Http response given a URL
func ReadResponse(u *url.URL) (*http.Response, error) {
	return readResponse(u)
}
