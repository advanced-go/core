package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	uri2 "github.com/advanced-go/core/uri"
	"net/http"
	"net/url"
	"os"
)

const (
	statusLoc         = PkgPath + ":ReadStatus"
	StatusOKUri       = "urn:status:ok"
	StatusNotFoundUri = "urn:status:notfound"
	StatusTimeoutUri  = "urn:status:timeout"
)

func ReadStatus(uri string) runtime.Status {
	status1 := constStatus(uri)
	if status1 != nil {
		return status1
	}
	/*
		if s, ok := t.(string); ok {
			//if len(s) == 0 || s == StatusOKUri {
			//	return runtime.StatusOK()
			//}
			uri = s
		} else {
			return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, errors.New(fmt.Sprintf("error: URI parameter is an invalid type: %v", reflect.TypeOf(t))))
		else {
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
	*/

	u, err := url.Parse(uri)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, statusLoc, err)
	}
	buf, err1 := os.ReadFile(uri2.FileName(u))
	if err1 != nil {
		return runtime.NewStatusError(runtime.StatusIOError, statusLoc, err1)
	}
	var status runtime.SerializedStatusState
	err = json.Unmarshal(buf, &status)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, statusLoc, err)
	}
	if len(status.Err) > 0 {
		return runtime.NewStatusError(status.Code, status.Location, errors.New(status.Err))
	}
	return runtime.NewStatus(status.Code).AddLocation(status.Location)
}

func constStatus(url string) runtime.Status {
	if len(url) == 0 {
		return runtime.StatusOK()
	}
	switch url {
	case StatusOKUri:
		return runtime.StatusOK()
	case StatusNotFoundUri:
		return runtime.NewStatus(http.StatusNotFound)
	case StatusTimeoutUri:
		return runtime.NewStatus(http.StatusGatewayTimeout)
	}
	return nil
}
