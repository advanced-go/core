package exchange

import (
	"errors"
	"net/http"
)

type Exchange func(req *http.Request) (*http.Response, error)
type Select func(req *http.Request) bool

type Proxy struct {
	Select Select
	Do     Exchange
}

var (
	list []Proxy
)

func AddProxy(p Proxy) error {
	if p.Select == nil || p.Do == nil {
		return errors.New("error: proxy Select or Do fn is nil")
	}
	list = append(list, p)
	return nil
}

// ProxyLookup - determine if a request maps to an Exchange function.
func ProxyLookup(req *http.Request) Exchange {
	if list != nil {
		for _, r := range list {
			if r.Select(req) {
				return r.Do
			}
		}
	}
	return nil
}
