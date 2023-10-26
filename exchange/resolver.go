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

type Resolver func(string) string

var resolver Resolver = defaultResolver

func SetResolver(f Resolver) {
	if f != nil {
		resolver = f
	}
}

// Resolve - given a string, resolve the string to an url.
func Resolve(s string) string {
	return resolver(s)
}

func defaultResolver(u string) string {
	// if an endpoint, then default to defaultOrigin
	if strings.HasPrefix(u, "/") {
		return defaultOrigin + u
	}
	// else pass through
	return u
}
