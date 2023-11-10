package log

import (
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// Wrap - wrap a DoHandler with access logging
func Wrap(handler runtime.DoHandler) runtime.DoHandler {
	return func(ctx any, req *http.Request, body any) (any, *runtime.Status) {
		var start = time.Now().UTC()

		if handler == nil {
			return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/Apply", errors.New("error:Do handler function is nil for access log")).SetRequestId(req.Context())
		}
		t, status := handler(ctx, req, body)
		InternalAccess(start, time.Since(start), req, &http.Response{StatusCode: status.Code()}, -1, "")
		return t, status
	}
}
