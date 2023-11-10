package http2

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, location string) {
	var e E

	if status.Content() == nil {
		return
	}
	buf, rc, status1 := WriteBytes(status.Content(), status.Header().Get(ContentType))
	if !status1.OK() {
		e.Handle(status, status.RequestId(), location+"/writeStatusContent")
		return
	}
	w.Header().Set(ContentType, rc)
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, location+"/writeStatusContent", err), "", "")
	}
}
