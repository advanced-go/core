package runtime

import (
	"compress/gzip"
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

	encodingReaderLoc = PkgPath + ":EncodingReader"
	encodingWriterLoc = PkgPath + ":EncodingWriter"
	encodingErrorFmt  = "error: content encoding not supported [%v]"
)

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

func EncodingReader(r io.Reader, h http.Header) (io.Reader, Status) {
	var reader io.Reader
	var err error

	encoding := contentEncoding(h)
	switch encoding {
	case GzipEncoding:
		reader, err = gzip.NewReader(r)
	case BrotliEncoding:
		err = encodingError(encoding)
	case DeflateEncoding:
		err = encodingError(encoding)
	case CompressEncoding:
		err = encodingError(encoding)
	case NoneEncoding:
		return r, StatusOK()
	default:
		err = encodingError(encoding)
	}
	if err != nil || reader == nil {
		return nil, NewStatusError(StatusContentDecodingError, encodingReaderLoc, err)
	}
	return reader, StatusOK()
}

func EncodingWriter(w io.Writer, h http.Header) (*gzip.Writer, Status) {
	var writer io.Writer
	var err error

	encoding := acceptEncoding(h)
	switch encoding {
	case GzipEncoding:
		return gzip.NewWriter(w), StatusOK()
	case BrotliEncoding:
		err = encodingError(encoding)
	case DeflateEncoding:
		err = encodingError(encoding)
	case CompressEncoding:
		err = encodingError(encoding)
	case NoneEncoding:
		return nil, StatusOK()
	default:
		err = encodingError(encoding)
	}
	if err != nil || writer == nil {
		return nil, NewStatusError(StatusContentDecodingError, encodingWriterLoc, err)
	}
	return nil, StatusOK()
}
