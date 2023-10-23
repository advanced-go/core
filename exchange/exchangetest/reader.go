package exchangetest

import (
	"io"
)

// ReaderCloser - test type for a body.ReadCloser interface
type readerCloser struct {
	Reader io.Reader
	Err    error
}

func (r *readerCloser) Read(p []byte) (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	return r.Reader.Read(p)
}

func (r *readerCloser) Close() error {
	return nil
}

func newReaderCloser(reader io.Reader, err error) *readerCloser {
	rc := new(readerCloser)
	rc.Reader = reader
	rc.Err = err
	return rc
}
