package messaging

import (
	"net/url"
	"strings"
)

const (
	urnSeparator = ":"
)

// uprootUrn - uproot an embedded urn in a uri
func uprootUrn(uri string) (nid, nss string, ok bool) {
	if uri == "" {
		return
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err.Error(), "", false
	}
	var str []string
	if u.Path[0] == '/' {
		str = strings.Split(u.Path[1:], urnSeparator)
	} else {
		str = strings.Split(u.Path, urnSeparator)
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
