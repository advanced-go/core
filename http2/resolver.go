package http2

import (
	"github.com/go-ai-agent/core/runtime"
	"strings"
)

type Resolver func(string) string

var (
	defaultOrigin          = "http://localhost:8080"
	resolver      Resolver = defaultResolver
	list          []Resolver
)

func SetDefaultOrigin(s string) {
	if !runtime.IsDebugEnvironment() {
		return
	}
	if len(s) != 0 {
		defaultOrigin = s
	}
}

func AddResolver(fn Resolver) {
	if !runtime.IsDebugEnvironment() || fn == nil {
		return
	}
	// do not need mutex, as this is only called from test
	list = append(list, fn)
}

// Resolve - resolve a string to an url.
func Resolve(s string) string {
	if !runtime.IsDebugEnvironment() {
		return ""
	}
	if list != nil {
		for _, r := range list {
			url := r(s)
			if len(url) != 0 {
				return url
			}
		}
	}
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
