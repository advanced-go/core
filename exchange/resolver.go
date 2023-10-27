package exchange

import (
	"errors"
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
	if len(s) != 0 {
		defaultOrigin = s
	}
}

func AddResolver(fn Resolver) error {
	if !runtime.IsDebugEnvironment() {
		return errors.New("error: adding a resolver is only available in DEBUG environment")
	}
	if fn == nil {
		return errors.New("error: resolver fn is nil")
	}
	// TODO : need to ensure mutex
	list = append(list, fn)
	return nil
}

// Resolve - resolve a string to an url.
func Resolve(s string) string {
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
