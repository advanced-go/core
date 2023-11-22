package urn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	Scheme = "urn"
)

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

func Build(nsid, nss string) string {
	return fmt.Sprintf("urn:%v:%v", nsid, nss)
}

func ToUri(scheme, host, urn string) string {
	i := len("urn:")
	return fmt.Sprintf("%v://%v/%v", scheme, host, urn[i:])
}

func FromUri(uri string) string {
	u, _ := url.Parse(uri)
	return fmt.Sprintf("urn:%v", u.Path[1:])

}
