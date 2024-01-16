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
		return ReadFile(ptr)
	case *url.URL:
		return ReadFile(ptr.String())
	case []byte:
		return ptr, StatusOK()
	case io.Reader:
		return ReadAll(io.NopCloser(ptr))
	case io.ReadCloser:
		return ReadAll(ptr)
	case *http.Response:
		return ReadAll(ptr.Body)
	case *http.Request:
		return ReadAll(ptr.Body)
	default:
	}
	return nil, NewStatusError(StatusInvalidArgument, newBytesLoc, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
}

// ReadFile - read a file
func ReadFile(uri string) ([]byte, Status) {
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

// ReadAll - read the body with a Status
func ReadAll(body io.ReadCloser) ([]byte, Status) {
	if body == nil {
		return nil, StatusOK()
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Printf("%v", err)
		}
	}(body)
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, NewStatusError(StatusIOError, readAllLoc, err)
	}
	return buf, StatusOK()
}
