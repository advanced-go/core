package exchange

import (
	"errors"
	"github.com/go-sre/core/runtime"
	"net/http"
	"strings"
)

const (
	SchemeName = "scheme"
	HostName   = "host"
	PathName   = "path"
	QueryName  = "query"
	MethodName = "method"
)

func LookupRequest(name string, req *http.Request) (string, error) {
	if req == nil {
		return "", errors.New("invalid argument: Request is nil")
	}
	switch strings.ToLower(name) {
	case SchemeName:
		return req.URL.Scheme, nil
	case HostName:
		return req.URL.Host, nil
	case PathName:
		return req.URL.Path, nil
	case QueryName:
		return req.URL.RawQuery, nil
	case MethodName:
		return req.Method, nil
	}
	return runtime.LookupEnv(name)
}
