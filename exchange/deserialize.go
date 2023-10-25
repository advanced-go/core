package exchange

import (
	"encoding/json"
	"errors"
	"github.com/go-ai-agent/core/httpx"
	"github.com/go-ai-agent/core/runtime"
	"io"
)

var deserializeLoc = PkgUri + "/Deserialize"

// Deserialize - provide deserialization of a request/response body
func Deserialize[T any](body io.ReadCloser) (T, *runtime.Status) {
	var t T

	if body == nil {
		return t, runtime.NewStatusError(errors.New("body is nil")).SetLocation(deserializeLoc).SetCode(runtime.StatusInvalidContent)
	}
	switch ptr := any(&t).(type) {
	case *[]byte:
		buf, status := httpx.ReadAll(body)
		if !status.OK() {
			return t, status
		}
		*ptr = buf
	default:
		err := json.NewDecoder(body).Decode(&t)
		if err != nil {
			return t, runtime.NewStatusError(err).SetLocation(deserializeLoc).SetCode(runtime.StatusJsonDecodeError)
		}
	}
	return t, runtime.NewStatusOK()
}
