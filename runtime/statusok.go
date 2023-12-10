package runtime

import (
	"fmt"
	"net/http"
	"time"
)

var statusOK = new(statusEmpty)

// StatusOK - return the single statusOK
func StatusOK() Status {
	return statusOK
}

type statusEmpty struct{}

func (s *statusEmpty) Code() int      { return http.StatusOK }
func (s *statusEmpty) OK() bool       { return true }
func (s *statusEmpty) NotFound() bool { return false }
func (s *statusEmpty) Http() int      { return http.StatusOK }

func (s *statusEmpty) IsErrors() bool    { return false }
func (s *statusEmpty) Errors() []error   { return nil }
func (s *statusEmpty) FirstError() error { return nil }

func (s *statusEmpty) Duration() time.Duration { return 0 }
func (s *statusEmpty) SetDuration(_ time.Duration) Status {
	return notImplementedSet("SetDuration()", s)
}

func (s *statusEmpty) RequestId() string         { return "" }
func (s *statusEmpty) SetRequestId(_ any) Status { return notImplementedSet("SetRequestId()", s) }

func (s *statusEmpty) Location() []string { return nil }

// AddLocation - allowed
func (s *statusEmpty) AddLocation(_ string) Status {
	return s //notImplementedSet("AddLocation()", s)
}

func (s *statusEmpty) IsContent() bool                 { return false }
func (s *statusEmpty) Content() any                    { return nil }
func (s *statusEmpty) ContentHeader() http.Header      { return nil }
func (s *statusEmpty) ContentString() string           { return "" }
func (s *statusEmpty) SetContent(_ any, _ bool) Status { return notImplementedSet("SetContent()", s) }

func (s *statusEmpty) Description() string { return "OK" }
func (s *statusEmpty) String() string      { return s.Description() }

func notImplementedSet(fn string, s Status) Status {
	fmt.Printf("function StatusOK.%v is not implemented\n", fn)
	return s
}
