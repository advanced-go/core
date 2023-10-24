package runtime

import "encoding/json"

func MarshalType[E ErrorHandler, T any](t T) ([]byte, *Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		var e E
		return nil, e.Handle(nil, "MarshalType", err)
	}
	return buf, NewStatusOK()
}

func UnmarshalType[E ErrorHandler, T any](buf []byte) (T, *Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		var e E
		return t, e.Handle(nil, "UnmarshalType", err)
	}
	return t, NewStatusOK()
}
