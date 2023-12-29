package http2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"reflect"
)

const (
	readAllLoc = PkgPath + ":ReadAll"
)

// ReadAll - read the type body with a runtime.Status
func ReadAll(t any) ([]byte, runtime.Status) {
	if t == nil {
		return nil, runtime.StatusOK()
	}
	if r, ok := t.(*http.Response); ok {
		return readAllBody(r.Body)
	}
	if body, ok := t.(io.ReadCloser); ok {
		return readAllBody(body)
	}
	return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, readAllLoc, errors.New(fmt.Sprintf("%v", reflect.TypeOf(t))))
}

// readAllBody - read the body with a runtime.Status
func readAllBody(body io.ReadCloser) ([]byte, runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
		}
	}(body)
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, readAllLoc, err)
	}
	return buf, runtime.StatusOK()
}
