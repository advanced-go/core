package runtime

import (
	"context"
	"encoding/json"
)

func MarshalType[E ErrorHandler, T any](ctx context.Context, t T) ([]byte, *Status) {
	buf, err := json.Marshal(t)
	if err != nil {
		var e E
		return nil, e.Handle(ctx, "MarshalType", err)
	}
	return buf, NewStatusOK()
}

func UnmarshalType[E ErrorHandler, T any](ctx context.Context, buf []byte) (T, *Status) {
	var t T

	err := json.Unmarshal(buf, &t)
	if err != nil {
		var e E
		return t, e.Handle(ctx, "UnmarshalType", err)
	}
	return t, NewStatusOK()
}
