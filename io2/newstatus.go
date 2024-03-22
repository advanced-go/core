package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"os"
)

const (
	//StatusOKUri       = "urn:status:ok"
	//StatusNotFoundUri = "urn:status:notfound"
	//StatusTimeoutUri  = "urn:status:timeout"
	//newLoc            = PkgPath + ":New"
	newStatusLoc = PkgPath + ":NewS"
)

type serializedStatusState struct {
	Code     int    `json:"code"`
	Location string `json:"location"`
	Err      string `json:"err"`
}

// NewStatusFrom - create a new Status from a URI
func NewStatusFrom(uri string) *runtime.Status {
	status := statusFromConst(uri)
	if status != nil {
		return status
	}
	status = ValidateUri(uri)
	if !status.OK() {
		return status
	}
	buf, err1 := os.ReadFile(FileName(uri))
	if err1 != nil {
		return runtime.NewStatusError(runtime.StatusIOError, err1)
	}
	var status2 serializedStatusState
	err := json.Unmarshal(buf, &status2)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
	}
	if len(status2.Err) > 0 {
		return runtime.NewStatusError(status2.Code, errors.New(status2.Err))
	}
	return runtime.NewStatus(status2.Code).AddLocation()
}

func statusFromConst(url string) *runtime.Status {
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
