package runtime

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	AcceptEncoding    = "Accept-Encoding"
	encodingReaderLoc = PkgPath + ":EncodingReader"
	encodingWriterLoc = PkgPath + ":EncodingWriter"
)

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

/*

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
		return nil, StatusOK()
	default:
		err = encodingError(encoding)
	}
	if err != nil {
		return nil, NewStatusError(StatusContentEncodingError, encodingReaderLoc, err)
	}
	return reader, StatusOK()
}


*/

func EncodingWriter(w io.Writer, h http.Header) (*gzip.Writer, Status) {
	//var writer io.Writer
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
	if err != nil {
		return nil, NewStatusError(StatusContentEncodingError, encodingWriterLoc, err)
	}
	return nil, StatusOK()
}
