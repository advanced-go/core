package controller

import (
	"encoding/json"
	"errors"
)

const (
	DefaultIngressRouteName = "default-ingress"
	DefaultEgressRouteName  = "default-egress"
)

// ReadRoutes - read routes from the []byte representation of a route configuration
func ReadRoutes(buf []byte) ([]Route, error) {
	var config []RouteConfig

	if buf == nil {
		return nil, errors.New("invalid argument: buffer is nil")
	}
	err1 := json.Unmarshal(buf, &config)
	if err1 != nil {
		return nil, err1
	}
	var routes []Route
	for _, c := range config {
		r, err := NewRouteFromConfig(c)
		if err != nil {
			return nil, err
		}
		routes = append(routes, r)
	}
	return routes, nil
}

// AddRoutes - read the routes from the []byte and create the CtrlTable controller entries
func AddRoutes(buf []byte) ([]Route, []error) {
	routes, err := ReadRoutes(buf)
	if err != nil {
		return routes, []error{err}
	}
	var errs []error
	for _, r := range routes {
		switch r.Name {
		default:
			errs = EgressTable().AddController(r)
		}
		if len(errs) > 0 {
			return nil, errs
		}
	}
	return routes, nil
}

func InitControllers(read func() ([]byte, error), update func(routes []Route) error) []error {
	if read == nil || update == nil {
		return []error{errors.New("invalid argument: read or updater function is nil")}
	}
	buf, err := read()
	if err != nil {
		return []error{err}
	}
	routes, errs := AddRoutes(buf)
	if len(errs) > 0 {
		return errs
	}
	err = update(routes)
	if err != nil {
		return []error{err}
	}
	return nil
}

func SetAction(name string, action Actuator) {
	EgressTable().SetAction(name, action)
}

func SetUriMatcher(fn UriMatcher) {
	EgressTable().SetUriMatcher(fn)
}

func SetHttpMatcher(fn HttpMatcher) {
	EgressTable().SetHttpMatcher(fn)
}
