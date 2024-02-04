package uri

import (
	"net/url"
	"strings"
)

const (
	UrnScheme    = "urn"
	UrnSeparator = ":"
)

// UprootUrn - uproot an embedded urn in a uri
func UprootUrn(uri string) (nid, nss string, ok bool) {
	if uri == "" {
		return
	}
	if strings.HasPrefix(uri, UrnScheme) {
		return uri, "", true
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
