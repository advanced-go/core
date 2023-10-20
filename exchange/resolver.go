package exchange

import "net/url"

var defaultOrigin = "http://localhost:8080"

func SetDefaultOrigin(s string) {
	if len(s) != 0 {
		defaultOrigin = s
	}
}

type UrlResolver func(string) string

var resolver UrlResolver = defaultResolver

func SetUrlResolver(f UrlResolver) {
	if f != nil {
		resolver = f
	}
}

func ResolveUrl(s string) string {
	return resolver(s)
}

func defaultResolver(u string) string {
	url, err := url.Parse(u)
	if err != nil {
		return err.Error()
	}
	if len(url.Host) == 0 {
		return defaultOrigin + u
	}
	return u
}
