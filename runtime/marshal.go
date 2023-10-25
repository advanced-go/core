package runtime

import (
	"encoding/json"
)

var marshalLoc = PkgUri + "/marshalType"

func MarshalType(t any) ([]byte, *Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		return nil, NewStatusError(err).SetLocation(PkgUri + "/marshalType")
	}
	return buf, NewStatusOK()
}

func UnmarshalType[T any](buf []byte) (T, *Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		return t, NewStatusError(err).SetLocation(PkgUri + "/unmarshalType")
	}
	return t, NewStatusOK()
}
