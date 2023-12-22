package uri

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	UrnScheme    = "urn"
	UrnSeparator = ":"
	FileScheme   = "file"
)

func IsFileScheme(u *url.URL) bool {
	if u == nil {
		return false
	}
	return u.Scheme == FileScheme
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

// UprootUrn - uproot an embedded urn in a uri
func UprootUrn(uri string) (nid, nss string, ok bool) {
	if uri == "" {
		return
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err.Error(), "", false
	}
	var str []string
	if u.Path[0] == '/' {
		str = strings.Split(u.Path[1:], UrnSeparator)
	} else {
		str = strings.Split(u.Path, UrnSeparator)
	}
	switch len(str) {
	case 0:
		return
	case 1:
		return str[0], "", true
	case 2:
		nid = str[0]
		nss = str[1]
	}
	return nid, nss, true
}

// Parse - urn safe Uri parser
func Parse(urn string) (scheme, ndis, nss string, err error) {
	if urn == "" {
		return
	}
	str := strings.Split(urn, ":")
	if len(str) < 3 {
		return "", "", "", errors.New(fmt.Sprintf("invalid urn format : %v", urn))
	}
	return str[0], str[1], str[2], nil
}

func Build(nid, nss string) string {
	return fmt.Sprintf("urn:%v:%v", nid, nss)
}

func ToUri(scheme, host, urn string) string {
	i := len("urn:")
	return fmt.Sprintf("%v://%v/%v", scheme, host, urn[i:])
}

func FromUri(uri string) string {
	u, _ := url.Parse(uri)
	return fmt.Sprintf("urn:%v", u.Path[1:])

}
