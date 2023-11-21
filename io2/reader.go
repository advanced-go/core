package io2

import (
	"github.com/advanced-go/core/runtime"
	"io"
)

// ReadAll - read the body with a runtime.Status
func ReadAll(body io.ReadCloser) ([]byte, runtime.Status) {
	if body == nil {
		return nil, runtime.NewStatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, PkgPath+"/ReadAll", err)
	}
	return buf, runtime.NewStatusOK()
}
