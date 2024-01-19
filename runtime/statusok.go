package runtime

import (
	"fmt"
	"net/http"
	"time"
)

var okStatus = new(statusOK)

// StatusOK - return the single status OK
func StatusOK() Status {
	return okStatus
}

type statusOK struct{}

func (s *statusOK) Code() int      { return http.StatusOK }
func (s *statusOK) OK() bool       { return true }
func (s *statusOK) NotFound() bool { return false }
func (s *statusOK) Http() int      { return http.StatusOK }

func (s *statusOK) IsErrors() bool     { return false }
func (s *statusOK) ErrorList() []error { return nil }
func (s *statusOK) Error() error       { return nil }

func (s *statusOK) Duration() time.Duration { return 0 }
func (s *statusOK) SetDuration(_ time.Duration) Status {
	return notImplementedSet("SetDuration()", s)
}

func (s *statusOK) RequestId() string         { return "" }
func (s *statusOK) SetRequestId(_ any) Status { return notImplementedSet("SetRequestId()", s) }

func (s *statusOK) Location() []string { return nil }

// AddLocation - allowed
func (s *statusOK) AddLocation(_ string) Status {
	return s //notImplementedSet("AddLocation()", s)
}

func (s *statusOK) Description() string { return "OK" }
func (s *statusOK) String() string      { return s.Description() }

func notImplementedSet(fn string, s Status) Status {
	fmt.Printf("function StatusOK.%v is not implemented\n", fn)
	return s
}
