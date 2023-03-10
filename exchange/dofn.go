package exchange

import (
	"github.com/gotemplates/core/runtime"
	"net/http"
)

var doLocation = pkgPath + "/dot"

func DoT[E runtime.ErrorHandler, T any, H HttpExchange](req *http.Request) (resp *http.Response, t T, status *runtime.Status) {
	var e E
	var h H

	if req == nil {
		return nil, t, runtime.NewHttpStatusCode(http.StatusInternalServerError)
	}
	var err error
	resp, err = h.Do(req)
	if err != nil {
		return nil, t, e.HandleWithContext(req.Context(), doLocation, err)
	}
	if resp == nil {
		return nil, t, runtime.NewHttpStatusCode(http.StatusInternalServerError)
	}
	if resp.StatusCode != http.StatusOK {
		return resp, t, runtime.NewHttpStatusCode(resp.StatusCode)
	}
	t, status = Deserialize[E, T](resp.Body)
	return
}
