package controller

// Route - route data
type Route struct {
	Name        string
	Pattern     string
	Timeout     *TimeoutConfig
	RateLimiter *RateLimiterConfig
}

type TimeoutConfigJson struct {
	Enabled    bool
	StatusCode int
	Duration   string
}

type RouteConfig struct {
	Name        string
	Pattern     string
	Traffic     string // Egress/Ingress
	Ping        bool   // Health traffic
	Protocol    string // gRPC, HTTP10, HTTP11, HTTP2, HTTP3
	Timeout     *TimeoutConfigJson
	RateLimiter *RateLimiterConfig
}

func newRoute(name string, config ...any) Route {
	return NewRoute(name, "", config...)
}

// NewRoute - creates a new route
func NewRoute(name string, protocol string, config ...any) Route {
	route := Route{}
	route.Name = name
	for _, cfg := range config {
		if cfg == nil {
			continue
		}
		switch c := cfg.(type) {
		case *TimeoutConfig:
			route.Timeout = c
		case *RateLimiterConfig:
			route.RateLimiter = c
		}
	}
	return route
}

// NewRouteFromConfig - creates a new route from configuration
func NewRouteFromConfig(config RouteConfig) (Route, error) {
	route := Route{}
	route.Name = config.Name
	route.Pattern = config.Pattern
	route.RateLimiter = config.RateLimiter
	if config.Timeout != nil {
		duration, err := ParseDuration(config.Timeout.Duration)
		if err != nil {
			return Route{}, err
		}
		route.Timeout = NewTimeoutConfig(config.Timeout.Enabled, config.Timeout.StatusCode, duration)
	}
	return route, nil
}

func (r Route) IsConfigured() bool {
	return r.Timeout != nil || r.RateLimiter != nil
}
