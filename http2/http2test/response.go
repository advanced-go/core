package http2test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	readResponseLocation = ":readResponse"
	fileExistsError      = "The system cannot find the file specified"
	fileScheme           = "file"
)

// readResponse - read a Http response given a URL
func readResponse(u *url.URL) (*http.Response, *runtime.Status) {
	serverErr := &http.Response{StatusCode: http.StatusInternalServerError, Status: "Internal Error"}

	if u == nil {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New("error: URL is nil"))
	}
	if u.Scheme != fileScheme {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New(fmt.Sprintf("error: Invalid URL scheme : %v", u.Scheme)))
	}
	buf, err := os.ReadFile(io2.FileName(u))
	if err != nil {
		if strings.Contains(err.Error(), fileExistsError) {
			return &http.Response{StatusCode: http.StatusNotFound, Status: "Not Found"}, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, err)
		}
		return serverErr, runtime.NewStatusError(runtime.StatusIOError, readResponseLocation, err)
	}
	resp1, err2 := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	if err2 != nil {
		return serverErr, runtime.NewStatusError(runtime.StatusIOError, readResponseLocation, err2)
	}
	return resp1, runtime.StatusOK()

}
