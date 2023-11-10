package httpxtest

import (
	"bufio"
	"bytes"
	"io"
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
