package exchange

import (
	"errors"
	"net/http"
)

type Exchange func(req *http.Request) (*http.Response, error)
type Valid func(req *http.Request) bool

type Resolver struct {
	Valid Valid
	Do    Exchange
}

var (
	list []Resolver
)

func AddResolver(r Resolver) error {
	if r.Valid == nil || r.Do == nil {
		return errors.New("error: Valid fn or Do fn is nil")
	}
	list = append(list, r)
	return nil
}

// Resolve - resolve a request to an Exchange function.
func Resolve(req *http.Request) Exchange {
	if list != nil {
		for _, r := range list {
			if r.Valid(req) {
				return r.Do
			}
		}
	}
	return nil
}
