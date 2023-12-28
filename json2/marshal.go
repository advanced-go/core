package json2

import (
	"encoding/json"
	"github.com/advanced-go/core/runtime"
)

const (
	marshalLoc   = ":Marshal"
	unMarshalLoc = ":Unmarshal"
)

// Marshal - JSON marshal with runtime.Status
func Marshal(t any) ([]byte, runtime.Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusJsonEncodeError, PkgPath+marshalLoc, err)
	}
	return buf, runtime.StatusOK()
}

// Unmarshal - JSON unmarshal with runtime.Status
func Unmarshal(buf []byte, t any) runtime.Status {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, PkgPath+unMarshalLoc, err)
	}
	return runtime.StatusOK()
}
