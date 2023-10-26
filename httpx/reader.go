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
		return nil, runtime.NewStatusError(runtime.StatusIOError, pkgUri+"/ReadAll", err)
	}
	return buf, runtime.NewStatusOK()
}

// ReaderCloser - test type for a body.ReadCloser interface
type ReaderCloser struct {
	Reader io.Reader
	Err    error
}

func (r *ReaderCloser) Read(p []byte) (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	return r.Reader.Read(p)
}

func (r *ReaderCloser) Close() error {
	return nil
}

func NewReaderCloser(reader io.Reader, err error) *ReaderCloser {
	rc := new(ReaderCloser)
	rc.Reader = reader
	rc.Err = err
	return rc
}
