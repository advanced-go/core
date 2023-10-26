package runtime

import (
	"encoding/json"
)

func MarshalType(t any) ([]byte, *Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, NewStatusError(StatusJsonEncodeError, pkgUri+"/MarshalType", err)
	}
	return buf, NewStatusOK()
}

func UnmarshalType[T any](buf []byte) (T, *Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		return t, NewStatusError(StatusJsonDecodeError, pkgUri+"/UnmarshalType", err)
	}
	return t, NewStatusOK()
}
