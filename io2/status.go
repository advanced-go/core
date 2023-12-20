package io2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	uri2 "github.com/advanced-go/core/uri"
	"net/url"
	"os"
	"reflect"
)

const (
	statusLoc   = ":ReadStatus"
	StatusOKUri = "urn:status:ok"
)

type statusState2 struct {
	Code     int    `json:"code"`
	Location string `json:"location"`
	Err      string `json:"err"`
}

func ReadStatus(t any) runtime.Status {
	uri := ""

	if t == nil {
		return runtime.StatusOK()
	}
	if s, ok := t.(string); ok {
		if len(s) == 0 || s == StatusOKUri {
			return runtime.StatusOK()
		}
		uri = s
	} else {
		if l, ok1 := t.([]string); ok1 {
			if len(l) == 0 || len(l) == 1 {
				return runtime.StatusOK()
			}
			if len(l[1]) == 0 || l[1] == StatusOKUri {
				return runtime.StatusOK()
			}
			uri = l[1]
		} else {
			return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, errors.New(fmt.Sprintf("error: URI parameter is an invalid type: %v", reflect.TypeOf(t))))
		}
	}
	//if len(uri) == 0 || uri == StatusOKUri {
	//	return runtime.StatusOK()
	//}
	u, err := url.Parse(uri)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, err)
	}
	buf, err1 := os.ReadFile(uri2.FileName(u))
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
