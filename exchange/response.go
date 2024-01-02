package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"io"
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
func readResponse(u *url.URL) (*http.Response, runtime.Status) {
	serverErr := &http.Response{StatusCode: http.StatusInternalServerError, Status: internalError}

	if u == nil {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New("error: URL is nil"))
	}
	if !uri.IsFileScheme(u) {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New(fmt.Sprintf("error: Invalid URL scheme : %v", u.Scheme)))
	}
	buf, err := os.ReadFile(uri.FileName(u))
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

type responseWriter struct {
	statusCode int
	header     http.Header
	bytes      bytes.Buffer
}

func newResponseWriter() *responseWriter {
	return new(responseWriter)
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return w.bytes.Write(p)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseWriter) Result() *http.Response {
	r := new(http.Response)
	if w.statusCode == 0 {
		r.StatusCode = http.StatusOK
	} else {
		r.StatusCode = w.statusCode
	}
	r.Header = w.header
	r.Body = io.NopCloser(bytes.NewReader(w.bytes.Bytes()))
	return r
}
