package http2test

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/advanced-go/core/http2"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"net/url"
)

// ParseRaw - parse a raw Uri without error
func ParseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

func ReadContent(rawHttp []byte) (*bytes.Buffer, error) {
	var content = new(bytes.Buffer)
	var writeTo = false

	reader := bufio.NewReader(bytes.NewReader(rawHttp))
	for {
		line, err := reader.ReadString('\n')
		if len(line) <= 2 && !writeTo {
			writeTo = true
			continue
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				return nil, err
			}
		}
		if writeTo {
			//fmt.Printf("%v", line)
			content.Write([]byte(line))
		}
	}
	return content, nil
}

func ReadContentFromLocation(h http.Header) ([]byte, runtime.Status) {
	content := h.Get(http2.ContentLocation)
	if len(content) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", errors.New("error: content-location is empty"))
	}
	u, _ := url.Parse(content)
	buf, err := io2.ReadFile(u)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", err)
	}
	return buf, runtime.NewStatusOK()
}
