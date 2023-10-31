package json

import (
	"encoding/json"
	"github.com/go-ai-agent/core/runtime"
)

func Marshal(t any) ([]byte, *runtime.Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusJsonEncodeError, runtime.PkgUri+"/Marshal", err)
	}
	return buf, runtime.NewStatusOK()
}

func Unmarshal(buf []byte, t any) *runtime.Status {
	err := json.Unmarshal(buf, t)
	if err != nil {
		return runtime.NewStatusError(runtime.StatusJsonDecodeError, runtime.PkgUri+"/Unmarshal", err)
	}
	return runtime.NewStatusOK()
}

/*
func Unmarshal[T any](buf []byte) (T, *runtime.Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, runtime.PkgUri+"/UnmarshalType", err)
	}
	return t, runtime.NewStatusOK()
}


*/
