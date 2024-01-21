package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	writeLoc = PkgPath + ":WriteResponse"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Content types supported: []byte, string, error, io.Reader, io.ReadCloser. Other types will be treated as JSON and serialized, if
// the headers content type is JSON. If not JSON, then an error will be raised.
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, content any, status runtime.Status, headers any) {
	var e E

	if status == nil {
		status = runtime.StatusOK()
	}
	SetHeaders(w, headers)
	//ct := w.Header().Get(ContentType)
	//accept := w.Header().Get(AcceptEncoding)
	//w.Header().Del(AcceptEncoding)
	w.WriteHeader(status.Http())
	_, status0 := writeContent(w, content, w.Header().Get(ContentType))
	if !status0.OK() {
		e.Handle(status0, status0.RequestId(), writeLoc)
	}
	return
}

/*
	ct := GetContentType(headers)
	buf, status0 := runtime.Bytes(content, ct)
	if !status0.OK() {
		e.Handle(status0, status.RequestId(), writeLoc)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SetHeaders(w, headers)
	if len(ct) == 0 {
		w.Header().Set(ContentType, http.DetectContentType(buf))
	}
	w.WriteHeader(status.Http())

*/
/*
	bytes, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, err), status.RequestId(), "")
	}
	if bytes != len(buf) {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, errors.New(fmt.Sprintf("error on ResponseWriter().Write() -> [got:%v] [want:%v]\n", bytes, len(buf)))), status.RequestId(), "")
	}

*/

/*
func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status runtime.Status, location string) {
	var e E

	if status.Content() == nil {
		return
	}
	ct := status.ContentHeader().Get(ContentType)
	buf, status1 := runtime.Bytes(status.Content(), ct)
	if !status1.OK() {
		e.Handle(status, status.RequestId(), location+writeStatusContentLoc)
		return
	}
	if len(ct) == 0 {
		ct = http.DetectContentType(buf)
	}
	w.Header().Set(ContentType, ct)
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, location+writeStatusContentLoc, err), "", "")
	}
}


*/
