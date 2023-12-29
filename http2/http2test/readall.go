package http2test

import (
	"github.com/advanced-go/core/runtime"
	"io"
)

const (
	PkgPath = "github.com/advanced-go/http2/http2test"
)

// ReadAll - read the body with a runtime.Status
func ReadAll(body io.ReadCloser) ([]byte, runtime.Status) {
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
		return nil, runtime.NewStatusError(runtime.StatusIOError, PkgPath+":ReadAll", err)
	}
	return buf, runtime.StatusOK()
}
