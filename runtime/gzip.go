package runtime

import (
	"compress/gzip"
	"io"
)

const (
	gzipWriterLoc = PkgPath + ":GzipWriter.Close()"
	gzipReaderLoc = PkgPath + ":GzipReader.Init()"
)

type gzipWriter struct {
	writer *gzip.Writer
}

func NewGzipWriter(w io.Writer) EncodingWriter {
	zw := new(gzipWriter)
	zw.writer = gzip.NewWriter(w)
	return zw
}

func (g *gzipWriter) Write(p []byte) (n int, err error) {
	return g.writer.Write(p)
}

func (g *gzipWriter) ContentEncoding() string {
	return GzipEncoding
}

func (g *gzipWriter) Close() Status {
	var errs []error

	err := g.writer.Flush()
	if err != nil {
		errs = append(errs, err)
	}
	err = g.writer.Close()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return NewStatusError(StatusGzipEncodingError, gzipWriterLoc, errs...)
	}
	return StatusOK()
}

type gzipReader struct {
	reader *gzip.Reader
}

func NewGzipReader(r io.Reader) (EncodingReader, Status) {
	zr := new(gzipReader)
	var err error
	zr.reader, err = gzip.NewReader(r)
	if err != nil {
		return nil, NewStatusError(StatusGzipEncodingError, gzipReaderLoc, err)
	}
	return zr, StatusOK()
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.reader.Read(p)
}

func (g *gzipReader) Close() Status {
	err := g.reader.Close()
	if err != nil {
		return NewStatusError(StatusGzipEncodingError, gzipReaderLoc, err)
	}
	return StatusOK()
}
