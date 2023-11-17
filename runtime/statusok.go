package runtime

import (
	"net/http"
	"time"
)

var status2 = newStatusOK()

func NewStatusOK() Status {
	return status2
}

type statusOK struct{}

func newStatusOK() Status     { return new(statusOK) }
func (s *statusOK) Code() int { return http.StatusOK }
func (s *statusOK) OK() bool  { return true }
func (s *statusOK) Http() int { return http.StatusOK }

func (s *statusOK) IsErrors() bool    { return false }
func (s *statusOK) Errors() []error   { return nil }
func (s *statusOK) FirstError() error { return nil }

func (s *statusOK) Duration() time.Duration            { return 0 }
func (s *statusOK) SetDuration(_ time.Duration) Status { return s }

func (s *statusOK) RequestId() string         { return "" }
func (s *statusOK) SetRequestId(_ any) Status { return s }

func (s *statusOK) Location() []string          { return nil }
func (s *statusOK) AddLocation(_ string) Status { return s }

func (s *statusOK) IsContent() bool                 { return false }
func (s *statusOK) Content() any                    { return nil }
func (s *statusOK) ContentHeader() http.Header      { return nil }
func (s *statusOK) ContentString() string           { return "" }
func (s *statusOK) SetContent(_ any, _ bool) Status { return s }

func (s *statusOK) Description() string { return "OK" }
func (s *statusOK) String() string      { return s.Description() }
