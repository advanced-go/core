package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	writeStatusContentLoc = ":writeStatusContent"
)

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status runtime.Status, location string) {
	var e E

	if status.Content() == nil {
		return
	}
	ct := status.ContentHeader().Get(ContentType)
	buf, status1 := WriteBytes(status.Content(), ct)
	if !status1.OK() {
		e.Handle(status, status.RequestId(), location+writeStatusContentLoc)
		return
	}
	if len(ct) == 0 {
		ct = http.DetectContentType(buf)
	}
	w.Header().Set(ContentType, ct)
	//w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, location+writeStatusContentLoc, err), "", "")
	}
}
