package exchangetest

import "bytes"

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
	contentTypeText = "text/plain"
	jsonSuffix      = ".json"
)

//
//path := u.Path
//if strings.HasPrefix(path, "/") {
//path = path[1:]
//}

//}
/*else {
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io2.NopCloser(bytes.NewReader(buf))}
	if strings.HasSuffix(path, jsonSuffix) {
		resp.Header.Add(contentType, contentTypeJson)
	} else {
		resp.Header.Add(contentType, contentTypeText)
	}
	return resp, runtime.StatusOK()
}
*/
//return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatus(http.StatusInternalServerError)
var (
	http11Bytes = []byte("HTTP/1.1")
	http12Bytes = []byte("HTTP/1.2")
	http20Bytes = []byte("HTTP/2.0")
)

func isHttpResponseMessage(buf []byte) bool {
	if buf == nil {
		return false
	}
	l := len(http11Bytes)
	if bytes.Equal(buf[0:l], http11Bytes) {
		return true
	}
	l = len(http12Bytes)
	if bytes.Equal(buf[0:l], http12Bytes) {
		return true
	}
	l = len(http20Bytes)
	if bytes.Equal(buf[0:l], http20Bytes) {
		return true
	}
	return false
}
