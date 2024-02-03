package io2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	readFileLoc = PkgPath + ":ReadFile"
	readAllLoc  = PkgPath + ":ReadAll"
)

// ReadFile - read a file with a Status
func ReadFile(uri string) ([]byte, *runtime.Status) {
	status := ValidateUri(uri)
	if !status.OK() {
		return nil, status
	}
	buf, err := os.ReadFile(FileName(uri))
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, readFileLoc, err)
	}
	return buf, runtime.StatusOK()
}

// ReadAll - read the body with a Status
func ReadAll(body io.Reader, h http.Header) ([]byte, *runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
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
		return nil, runtime.NewStatusError(runtime.StatusIOError, readAllLoc, err)
	}
	return buf, runtime.StatusOK()
}

func ValidateUri(uri string) *runtime.Status {
	if len(uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, newStatusLoc, errors.New("error: URI is empty"))
	}
	if !strings.HasPrefix(uri, fileScheme) {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, newStatusLoc, errors.New(fmt.Sprintf("error: URI is not of scheme file: %v", uri)))
	}
	if !isJsonURL(uri) {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, newStatusLoc, errors.New("error: URI is not a JSON file"))
	}
	return runtime.StatusOK()
}
