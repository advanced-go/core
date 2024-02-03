package http2

import (
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	writeLoc     = PkgPath + ":WriteResponse"
	noneEncoding = ""
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Content types supported: []byte, string, error, io.Reader, io.ReadCloser. Other types will be treated as JSON and serialized, if
// the headers content type is JSON. If not JSON, then an error will be raised.
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, content any, status *runtime.Status, headers any) {
	var e E

	if status == nil {
		status = runtime.StatusOK()
	}
	SetHeaders(w, headers)
	if content == nil {
		w.WriteHeader(status.HttpCode())
		return
	}
	h := createAcceptEncoding(w.Header())
	writer, status0 := io2.NewEncodingWriter(w, h)
	if !status0.OK() {
		e.Handle(status0, runtime.RequestId(w.Header()), writeLoc)
		return
	}
	if writer.ContentEncoding() != io2.NoneEncoding {
		w.Header().Add(ContentEncoding, writer.ContentEncoding())
	}
	w.WriteHeader(status.HttpCode())
	_, status0 = writeContent(writer, content, w.Header().Get(ContentType))
	writer.Close()
	if !status0.OK() {
		e.Handle(status0, runtime.RequestId(w.Header()), writeLoc)
	}
}

func createAcceptEncoding(h http.Header) http.Header {
	out := make(http.Header)
	if h == nil {
		return out
	}
	accept := h.Get(AcceptEncoding)
	h.Del(AcceptEncoding)
	if len(accept) == 0 {
		return out
	}
	if len(h.Get(ContentEncoding)) != 0 {
		return out
	}
	out.Add(AcceptEncoding, accept)
	return out
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
