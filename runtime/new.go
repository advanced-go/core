package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if !uri2.IsFileScheme(uri) {
		return t, NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New(fmt.Sprintf("error: URI is not of scheme file: %v", uri)))
	}
	if !uri2.IsJson(uri) {
		return t, NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New("error: URI is not a JSON file"))
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
	status := statusFromConst(uri)
	if status != nil {
		return status
	}
	if !uri2.IsFileScheme(uri) {
		return NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New(fmt.Sprintf("error: URI is not of scheme file: %v", uri)))
	}
	if !uri2.IsJson(uri) {
		return NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New("error: URI is not a JSON file"))
	}
	buf, err1 := os.ReadFile(uri2.FileName(uri))
	if err1 != nil {
		return NewStatusError(StatusIOError, newStatusLoc, err1)
	}
	var status2 serializedStatusState
	err := json.Unmarshal(buf, &status2)
	if err != nil {
		return NewStatusError(StatusJsonDecodeError, newStatusLoc, err)
	}
	if len(status2.Err) > 0 {
		return NewStatusError(status2.Code, status2.Location, errors.New(status2.Err))
	}
	return NewStatus(status2.Code).AddLocation(status2.Location)
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
