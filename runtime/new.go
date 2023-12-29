package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	uri2 "github.com/advanced-go/core/uri"
	"io"
	"net/url"
	"os"
	"reflect"
)

const (
	StatusOKUri       = "urn:status:ok"
	StatusNotFoundUri = "urn:status:notfound"
	StatusTimeoutUri  = "urn:status:timeout"
	newStatusLoc      = PkgPath + ":NewS"
	newTypeLoc        = PkgPath + ":New"
)

// New - create a new type from JSON content
func New[T any](v any) (t T, status Status) {
	uri := ""
	switch ptr := any(v).(type) {
	case string:
		if uri2.IsStatusURL(ptr) {
			return t, NewS(ptr)
		}
		uri = ptr
	case *url.URL:
		if uri2.IsStatusURL(ptr.String()) {
			return t, NewS(ptr.String())
		}
		uri = ptr.String()
	case io.ReadCloser:
		err := json.NewDecoder(ptr).Decode(&t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newTypeLoc, err)
		}
		return t, StatusOK()
	default:
		return t, NewStatusError(StatusInvalidArgument, newTypeLoc, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
	}
	status = validateUri(uri)
	if !status.OK() {
		return
	}
	buf, err := os.ReadFile(uri2.FileName(uri))
	if err != nil {
		return t, NewStatusError(StatusIOError, newTypeLoc, err)
	}
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return t, NewStatusError(StatusJsonDecodeError, newTypeLoc, err)
	}
	return t, StatusOK()
}
