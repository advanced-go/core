package exchange

import "io"

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
