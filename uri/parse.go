package uri

import (
	"errors"
	"net/url"
	"strings"
)

const (
	FileScheme = "file"
)

// IsFileScheme - determine if a string, or URL uses the file scheme
func IsFileScheme(u any) bool {
	if u == nil {
		return false
	}
	if s, ok := u.(string); ok {
		return strings.HasPrefix(s, FileScheme)
	}
	if u2, ok := u.(*url.URL); ok {
		return u2.Scheme == FileScheme
	}
	return false
}

// ParseRaw - parse a raw Uri without error
func ParseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

// ParseUri - urn safe URI parser
func ParseUri(uri string) (scheme, host, path string) {
	if uri == "" {
		return
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err.Error(), "", ""
	}
	if u.Scheme == UrnScheme && u.Host == "" {
		t := strings.Split(u.Opaque, ":")
		if len(t) == 1 {
			return u.Scheme, t[0], ""
		}
		return u.Scheme, t[0], t[1]
	}
	return u.Scheme, u.Host, u.Path
}

// BuildUrl - build a URL expanding a template
func BuildUrl(url *url.URL, template string) (*url.URL, error) {
	if url == nil {
		return nil, errors.New("invalid parameter: URL is nil")
	}
	if template == "" {
		return url, nil
	}
	url2, err := Expand(template, func(name string) (string, error) {
		return LookupVariable(name, url)
	},
	)
	if err != nil {
		return nil, err
	}
	// Removing trailing "?" which happens if the template has a query variable, and the request URL does not
	// contain a query
	length := len(url2)
	if url2[length-1:] == "?" {
		url2 = url2[:length-1]
	}
	u, err1 := url.Parse(url2)
	if err1 != nil {
		return nil, err1
	}
	return u, nil
}
