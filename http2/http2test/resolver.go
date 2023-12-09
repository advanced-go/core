package http2test

import (
	"github.com/advanced-go/core/runtime"
	"strings"
)

type resolverFunc func(string) string

var (
	defaultOrigin = "http://localhost:8080"
	list          []resolverFunc
)

/*
	func SetDefaultOrigin(s string) {
		if !runtime.IsDebugEnvironment() {
			return
		}
		if len(s) != 0 {
			defaultOrigin = s
		}
	}
*/
func addResolver(fn resolverFunc) {
	if !runtime.IsDebugEnvironment() || fn == nil {
		return
	}
	// do not need mutex, as this is only called from test
	list = append(list, fn)
}

// resolve - resolve a string to an url.
func resolve(s string) string {
	if !runtime.IsDebugEnvironment() {
		return defaultResolver(s)
	}
	if list != nil {
		for _, r := range list {
			url := r(s)
			if len(url) != 0 {
				return url
			}
		}
	}
	return defaultResolver(s)
}

func defaultResolver(u string) string {
	// if an endpoint, then default to defaultOrigin
	if strings.HasPrefix(u, "/") {
		return defaultOrigin + u
	}
	// else pass through
	return u
}
