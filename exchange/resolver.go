package exchange

import (
	"strings"
)

var defaultOrigin = "http://localhost:8080"

func SetDefaultOrigin(s string) {
	if len(s) != 0 {
		defaultOrigin = s
	}
}

type UrlTextResolver func(string) string

var resolver UrlTextResolver = defaultResolver

func SetUrlResolver(f UrlTextResolver) {
	if f != nil {
		resolver = f
	}
}

func ResolveUrl(s string) string {
	return resolver(s)
}

func defaultResolver(u string) string {
	if strings.HasPrefix(u, "/") {
		return defaultOrigin + u
	}
	return u
}
