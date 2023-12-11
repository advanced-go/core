package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

const (
	statusLoc = PkgPath + ":ReadStatus"
	StatusOK  = "urn:status:ok"
)

type statusState2 struct {
	Code     int    `json:"code"`
	Location string `json:"location"`
	Err      string `json:"err"`
}

func ReadStatus(uri string) runtime.Status {
	if len(uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, errors.New("uir is empty"))
	}
	if uri == StatusOK {
		return runtime.StatusOK()
	}
	u, err := url.Parse(uri)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, err)
	}
	buf, err1 := ReadFile(u)
	if err1 != nil {
		return runtime.NewStatusError(runtime.StatusIOError, statusLoc, err1)
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
