package json2

import (
	"encoding/json"
	"github.com/advanced-go/core/runtime"
)

const (
	PkgPath = "github.com/advanced-go/core/json2"
)

// Marshal - JSON marshal with runtime.Status
func Marshal(t any) ([]byte, runtime.Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusJsonEncodeError, PkgPath+"/Marshal", err)
	}
	return buf, runtime.NewStatusOK()
}

// Unmarshal - JSON unmarshal with runtime.Status
func Unmarshal(buf []byte, t any) runtime.Status {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, PkgPath+"/Unmarshal", err)
	}
	return runtime.NewStatusOK()
}
