package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

const (
	statusLoc = PkgPath + ":ReadStatus"
)

type statusState2 struct {
	Code     int    `json:"code"`
	Location string `json:"location"`
	Err      string `json:"err"`
}

func ReadStatus(u *url.URL) runtime.Status {
	buf, err := ReadFile(u)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusIOError, statusLoc, err)
	}
	var status statusState2
	err = json.Unmarshal(buf, &status)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, statusLoc, err)
	}
	if len(status.Err) > 0 {
		return runtime.NewStatusError(status.Code, status.Location, errors.New(status.Err))
	}
	return runtime.NewStatus(status.Code).AddLocation(status.Location)
}
