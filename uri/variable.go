package uri

import (
	"errors"
	"net/url"
	"strings"
)

const (
	SchemeName = "scheme"
	HostName   = "host"
	PathName   = "path"
	QueryName  = "query"
	//MethodName = "method"
)

func LookupVariable(name string, url *url.URL) (string, error) {
	if url == nil {
		return "", errors.New("invalid argument: Url is nil")
	}
	switch strings.ToLower(name) {
	case SchemeName:
		return url.Scheme, nil
	case HostName:
		return url.Host, nil
	case PathName:
		return url.Path, nil
	case QueryName:
		return url.RawQuery, nil
		//case MethodName:
		//	return method, nil
	}
	return LookupEnv(name)
}
