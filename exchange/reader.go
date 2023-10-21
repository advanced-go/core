package exchange

import (
	"github.com/go-ai-agent/core/runtime"
	"io"
)

// ReadAll - read all the body, with a deferred close
func ReadAll[E runtime.ErrorHandler](body io.ReadCloser) ([]byte, *runtime.Status) {
	if body == nil {
		return nil, runtime.NewStatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		var e E
		return nil, e.Handle(nil, "ReadAll", err).SetCode(runtime.StatusIOError)
	}
	return buf, runtime.NewStatusOK()
}

// ReaderCloser - test type for a body.ReadCloser interface
type ReaderCloser struct {
	Reader io.Reader
}

func (r *ReaderCloser) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}

func (r *ReaderCloser) Close() error {
	return nil
}

func NewReaderCloser(reader io.Reader) *ReaderCloser {
	rc := new(ReaderCloser)
	rc.Reader = reader
	return rc
}
