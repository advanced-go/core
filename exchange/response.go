package exchange

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
	readResponseLocation = PkgPath + ":readResponse"
	fileExistsError      = "The system cannot find the file specified"
)

// readResponse - read a Http response given a URL
func readResponse(u *url.URL) (*http.Response, *runtime.Status) {
	serverErr := &http.Response{StatusCode: http.StatusInternalServerError, Status: internalError}

	if u == nil {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("error: URL is nil"), nil)
	}
	if u.Scheme != fileScheme {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("error: Invalid URL scheme : %v", u.Scheme)), nil)
	}
	buf, err := os.ReadFile(io2.FileName(u))
	if err != nil {
		if strings.Contains(err.Error(), fileExistsError) {
			return &http.Response{StatusCode: http.StatusNotFound, Status: "Not Found"}, runtime.NewStatusError(runtime.StatusInvalidArgument, err, nil)
		}
		return serverErr, runtime.NewStatusError(runtime.StatusIOError, err, nil)
	}
	resp1, err2 := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	if err2 != nil {
		return serverErr, runtime.NewStatusError(runtime.StatusIOError, err2, nil)
	}
	return resp1, runtime.StatusOK()

}
