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
	newStatusLoc = PkgPath + ":NewS"
)

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
	status = validateUri(uri)
	if !status.OK() {
		return status
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

func validateUri(uri string) Status {
	if len(uri) == 0 {
		return NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New("error: URI is empty"))
	}
	if !uri2.IsFileScheme(uri) {
		return NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New(fmt.Sprintf("error: URI is not of scheme file: %v", uri)))
	}
	if !uri2.IsJson(uri) {
		return NewStatusError(StatusInvalidArgument, newStatusLoc, errors.New("error: URI is not a JSON file"))
	}
	return StatusOK()
}
