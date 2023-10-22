package runtime

import (
	"errors"
	"net/url"
	"strings"
)

const (
	UrnScheme  = "urn"
	FileScheme = "file://"
)

// ParsePkgUrl - parse a package raw Uri
func ParsePkgUrl(rawUri string) *url.URL {
	u, err := url.Parse(FileScheme + rawUri)
	if err != nil {
		return nil
	}
	return u
}

// ParseRaw - parse a raw Uri without error
func ParseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

// ParseUri - urn safe Uri parser
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

func BuildUrl(url *url.URL, template string) (*url.URL, error) {
	if url == nil {
		return nil, errors.New("invalid parameter: URL is nil")
	}
	if template == "" {
		return url, nil
	}
	url2, err := Expand(func(name string) (string, error) {
		return LookupUrl(name, url)
	},
		template,
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
