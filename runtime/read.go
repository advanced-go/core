package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	ContentEncoding  = "Content-Encoding"
	GzipEncoding     = "gzip"
	BrotliEncoding   = "br"
	DeflateEncoding  = "deflate"
	CompressEncoding = "compress"
	NoneEncoding     = "none"

	encodingErrorFmt = "error: content encoding not supported [%v]"
	readFileLoc      = PkgPath + ":ReadFile"
	readAllLoc       = PkgPath + ":ReadAll"
)

// ReadFile - read a file with a Status
func ReadFile(uri string) ([]byte, Status) {
	status := validateUri(uri)
	if !status.OK() {
		return nil, status
	}
	buf, err := os.ReadFile(FileName(uri))
	if err != nil {
		return nil, NewStatusError(StatusIOError, readFileLoc, err)
	}
	return buf, StatusOK()
}

// ReadAll - read the body with a Status
func ReadAll(body io.Reader, h http.Header) ([]byte, Status) {
	if body == nil {
		return nil, StatusOK()
	}
	if rc, ok := any(body).(io.ReadCloser); ok {
		defer func() {
			err := rc.Close()
			if err != nil {
				fmt.Printf("error: io.ReadCloser.Close() [%v]", err)
			}
		}()
	}
	var buf []byte
	var err error

	encoding := contentEncoding(h)
	switch encoding {
	case GzipEncoding:
		zr, status := NewGzipReader(body)
		if !status.OK() {
			return nil, status.AddLocation(readAllLoc)
		}
		buf, err = io.ReadAll(zr)
		status = zr.Close()
		if !status.OK() {
			return nil, status.AddLocation(readAllLoc)
		}
	case NoneEncoding:
		buf, err = io.ReadAll(body)
	default:
		return nil, NewStatusError(StatusContentEncodingError, readAllLoc, encodingError(encoding))
	}
	if err != nil {
		return nil, NewStatusError(StatusIOError, readAllLoc, err)
	}
	return buf, StatusOK()
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
