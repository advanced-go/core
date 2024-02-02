package messaging

import (
	"fmt"
	"net/http"
	"time"
)

// Status - message status
type Status struct {
	Error    error
	Code     int
	Location string
	Duration time.Duration
}

func StatusOK() *Status {
	return NewStatus(http.StatusOK)
}

func NewStatus(code int) *Status {
	s := new(Status)
	s.Code = code
	return s
}

func NewStatusError(err error, location string) *Status {
	s := new(Status)
	s.Error = err
	s.Location = location
	return s
}

func NewStatusDuration(code int, duration time.Duration) *Status {
	s := new(Status)
	s.Code = code
	s.Duration = duration
	return s
}

func NewStatusDurationError(code int, duration time.Duration, err error) *Status {
	s := NewStatusDuration(code, duration)
	s.Error = err
	return s
}

func (s *Status) OK() bool {
	return s.Code == http.StatusOK
}

func (s *Status) String() string {
	if s.Error != nil {
		return fmt.Sprintf("%v", s.Error)
	} else {
		return fmt.Sprintf("%v", s.Code)
	}
}
