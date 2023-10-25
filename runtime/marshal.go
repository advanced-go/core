package runtime

import (
	"encoding/json"
)


func MarshalType(t any) ([]byte, *Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, NewStatusError(err).SetLocation(PkgUri + "/MarshalType")
	}
	return buf, NewStatusOK()
}

func UnmarshalType[T any](buf []byte) (T, *Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		return t, NewStatusError(err).SetLocation(PkgUri + "/UnmarshalType")
	}
	return t, NewStatusOK()
}
