package httptest

import "io"

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
