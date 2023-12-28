package exchangetest

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"net/url"
	"os"
)

const (
	readStatusLoc = ":readStatus"
)

func serverErrorResponse() *http.Response {
	resp := new(http.Response)
	resp.StatusCode = http.StatusInternalServerError
	resp.Status = "Internal Error"
	return resp
}

func testURL(req *http.Request) (*http.Response, runtime.Status) {
	if uri.IsStatusURL(req.URL.String()) {
		return readStatus(req.URL)
	}
	return nil, nil
}

func readStatus(u *url.URL) (*http.Response, runtime.Status) {
	if u != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusInvalidArgument, readStatusLoc, errors.New("invalid argument : URL is nil"))
	}
	buf, err1 := os.ReadFile(uri.FileName(u))
	if err1 != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusIOError, readStatusLoc, err1)
	}
	var status runtime.SerializedStatusState
	err := json.Unmarshal(buf, &status)
	if err != nil {
		return serverErrorResponse(), runtime.NewStatusError(runtime.StatusJsonDecodeError, readStatusLoc, err)
	}
	if len(status.Err) > 0 {
		return serverErrorResponse(), runtime.NewStatusError(status.Code, status.Location, errors.New(status.Err))
	}
	resp := new(http.Response)
	resp.StatusCode = status.Code
	return resp, runtime.NewStatus(status.Code).AddLocation(status.Location)
}
