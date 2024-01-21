package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	AcceptEncoding   = "Accept-Encoding"
	ContentEncoding  = "Content-Encoding"
	GzipEncoding     = "gzip"
	BrotliEncoding   = "br"
	DeflateEncoding  = "deflate"
	CompressEncoding = "compress"
	NoneEncoding     = "none"

	encodingErrorFmt  = "error: content encoding not supported [%v]"
	encodingReaderLoc = PkgPath + ":EncodingReader"
	encodingWriterLoc = PkgPath + ":EncodingWriter"
)

type EncodingReader interface {
	io.Reader
	Close() Status
}

type EncodingWriter interface {
	io.Writer
	Close() Status
}

func NewEncodingReader(r io.Reader, h http.Header) (EncodingReader, Status) {
	encoding := contentEncoding(h)
	switch encoding {
	case GzipEncoding:
		return NewGzipReader(r)
	case BrotliEncoding, DeflateEncoding, CompressEncoding:
		return nil, NewStatusError(StatusContentEncodingError, encodingReaderLoc, encodingError(encoding))
	default:
		return NewIdentityReader(r), StatusOK()
	}
}

func NewEncodingWriter(w io.Writer, h http.Header) (EncodingWriter, Status) {
	encoding := acceptEncoding(h)
	if strings.Contains(encoding, GzipEncoding) {
		return NewGzipWriter(w), StatusOK()
	}
	return NewIdentityWriter(w), StatusOK()
}

func encodingError(encoding string) error {
	return errors.New(fmt.Sprintf(encodingErrorFmt, encoding))
}

func contentEncoding(h http.Header) string {
	if h == nil {
		return NoneEncoding
	}
	enc := h.Get(ContentEncoding)
	if len(enc) > 0 {
		return strings.ToLower(enc)
	}
	return NoneEncoding
}

func acceptEncoding(h http.Header) string {
	if h == nil {
		return NoneEncoding
	}
	enc := h.Get(AcceptEncoding)
	if len(enc) > 0 {
		return strings.ToLower(enc)
	}
	return NoneEncoding
}

type identityReader struct {
	reader io.Reader
}

// NewIdentityReader - The default (identity) encoding; the use of no transformation whatsoever
func NewIdentityReader(r io.Reader) EncodingReader {
	ir := new(identityReader)
	ir.reader = r
	return ir
}

func (i *identityReader) Read(p []byte) (n int, err error) {
	return i.reader.Read(p)
}

func (i *identityReader) Close() Status {
	return StatusOK()
}

type identityWriter struct {
	writer io.Writer
}

// NewIdentityWriter - The default (identity) encoding; the use of no transformation whatsoever
func NewIdentityWriter(w io.Writer) EncodingWriter {
	iw := new(identityWriter)
	iw.writer = w
	return iw
}

func (i *identityWriter) Write(p []byte) (n int, err error) {
	return i.writer.Write(p)
}

func (i *identityWriter) Close() Status {
	return StatusOK()
}
