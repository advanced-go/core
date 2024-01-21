package runtime

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	readFileLoc = PkgPath + ":ReadFile"
	readAllLoc  = PkgPath + ":ReadAll"
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
	reader, status := NewEncodingReader(body, h)
	if !status.OK() {
		return nil, status.AddLocation(readAllLoc)
	}
	buf, err := io.ReadAll(reader)
	reader.Close()
	if err != nil {
		return nil, NewStatusError(StatusIOError, readAllLoc, err)
	}
	return buf, StatusOK()
}
