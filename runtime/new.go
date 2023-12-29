package runtime

import (
	"encoding/json"
	"errors"
	uri2 "github.com/advanced-go/core/uri"
	"net/http"
	"os"
)

const (
	StatusOKUri       = "urn:status:ok"
	StatusNotFoundUri = "urn:status:notfound"
	StatusTimeoutUri  = "urn:status:timeout"
	newStatusLoc      = PkgPath + ":NewS"
	newTypeLoc        = PkgPath + ":NewT"
)

func NewT[T any](uri string) (t T, status Status) {
	if len(uri) == 0 {
		return t, NewStatusError(StatusInvalidArgument, newTypeLoc, errors.New("error: URI is empty"))
	}
	buf, err := os.ReadFile(uri2.FileName(uri))
	if err != nil {
		return t, NewStatusError(StatusIOError, newTypeLoc, err)
	}
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return t, NewStatusError(StatusJsonDecodeError, newTypeLoc, err)
	}
	return t, StatusOK()
}

type serializedStatusState struct {
	Code     int    `json:"code"`
	Location string `json:"location"`
	Err      string `json:"err"`
}

func NewS(uri string) Status {
	status1 := statusFromConst(uri)
	if status1 != nil {
		return status1
	}
	buf, err1 := os.ReadFile(uri2.FileName(uri))
	if err1 != nil {
		return NewStatusError(StatusIOError, newStatusLoc, err1)
	}
	var status serializedStatusState
	err := json.Unmarshal(buf, &status)
	if err != nil {
		return NewStatusError(StatusJsonDecodeError, newStatusLoc, err)
	}
	if len(status.Err) > 0 {
		return NewStatusError(status.Code, status.Location, errors.New(status.Err))
	}
	return NewStatus(status.Code).AddLocation(status.Location)
}

func statusFromConst(url string) Status {
	if len(url) == 0 {
		return StatusOK()
	}
	switch url {
	case StatusOKUri:
		return StatusOK()
	case StatusNotFoundUri:
		return NewStatus(http.StatusNotFound)
	case StatusTimeoutUri:
		return NewStatus(http.StatusGatewayTimeout)
	}
	return nil
}

/*
func ReadStatus(uri string) runtime.Status {
	status1 := constStatus(uri)
	if status1 != nil {
		return status1
	}

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



*/
