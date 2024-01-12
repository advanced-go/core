package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

const (
	newBytesLoc = PkgPath + ":NewBytes"
	readAllLoc  = PkgPath + ":readAll"
)

// NewBytes - create a []byte from a type
func NewBytes(v any) ([]byte, Status) {
	switch ptr := v.(type) {
	case string:
		return readBytes(ptr)
	case *url.URL:
		return readBytes(ptr.String())
	case []byte:
		return ptr, StatusOK()
	case io.Reader:
		return readAll(io.NopCloser(ptr))
	case io.ReadCloser:
		return readAll(ptr)
	case *http.Response:
		return readAll(ptr.Body)
	default:
	}
	return nil, NewStatusError(StatusInvalidArgument, newBytesLoc, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
}

func readBytes(uri string) ([]byte, Status) {
	status := validateUri(uri)
	if !status.OK() {
		return nil, status
	}
	buf, err := os.ReadFile(FileName(uri))
	if err != nil {
		return nil, NewStatusError(StatusIOError, newBytesLoc, err)
	}
	return buf, StatusOK()
}

// readAll - read the body with a Status
func readAll(body io.ReadCloser) ([]byte, Status) {
	if body == nil {
		return nil, StatusOK()
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
		}
	}(body)
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, NewStatusError(StatusIOError, readAllLoc, err)
	}
	return buf, StatusOK()
}
