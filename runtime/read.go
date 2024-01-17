package runtime

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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

	readFileLoc       = PkgPath + ":ReadFile"
	readAllLoc        = PkgPath + ":ReadAll"
	encodingReaderLoc = PkgPath + ":EncodingReader"

	decodingErrorFmt = "error: content encoding not supported [%v]"
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
	reader, status := EncodingReader(body, h)
	if !status.OK() {
		return nil, status.AddLocation(readAllLoc)
	}
	buf, err1 := io.ReadAll(reader)
	if err1 != nil {
		return nil, NewStatusError(StatusIOError, readAllLoc, err1)
	}
	return buf, StatusOK()
}

func encodingError(encoding string) error {
	return errors.New(fmt.Sprintf(decodingErrorFmt, encoding))
}

func encodingHeader(h http.Header) string {
	if h == nil {
		return NoneEncoding
	}
	enc := h.Get(ContentEncoding)
	if len(enc) > 0 {
		return strings.ToLower(enc)
	}
	return NoneEncoding
}

func EncodingReader(r io.Reader, h http.Header) (io.Reader, Status) {
	var reader io.Reader
	var err error

	encoding := encodingHeader(h)
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
		//if rc, ok := any(r).(io.ReadCloser); ok {
		//	return rc, StatusOK()
		//}
		return r, StatusOK()
	default:
		err = encodingError(encoding)
	}
	if err != nil || reader == nil {
		return nil, NewStatusError(StatusContentDecodingError, encodingReaderLoc, err)
	}
	return reader, StatusOK()
}
