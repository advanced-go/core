package httpx

import (
	"github.com/go-ai-agent/core/runtime"
	"io"
)

// ReadAll - read all the body, with a deferred close
func ReadAll(body io.ReadCloser) ([]byte, *runtime.Status) {
	if body == nil {
		return nil, runtime.NewStatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, PkgUri+"/ReadAll", err)
	}
	return buf, runtime.NewStatusOK()
}
