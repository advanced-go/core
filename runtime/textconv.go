package runtime

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	comment   = "//"
	delimiter = ":"
)

type TextPair struct {
	Key   string
	Value string
}

// ValidateMap - validates a configuration map, iterating through all keys
func ValidateMap(m map[string]string, err error, keys ...string) (errs []error) {
	if m == nil {
		return []error{errors.New("config map is nil")}
	}
	if err != nil {
		errs = append(errs, errors.New(fmt.Sprintf("config map read error: %v", err)))
	}
	for _, k := range keys {
		if v, ok := m[k]; ok {
			if v == "" {
				errs = append(errs, errors.New(fmt.Sprintf("config map error: value for key does not exist [%v]", k)))
			}
		} else {
			errs = append(errs, errors.New(fmt.Sprintf("[config map error: key does not exist [%v]", k)))
		}
	}
	return
}

// TextToSlice - create a TextPari slice from a []byte
func TextToSlice(buf []byte) ([]TextPair, error) {
	var list []TextPair

	if len(buf) == 0 {
		return list, nil
	}
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		k, v, err0 := parseLine(line, false)
		if err0 != nil {
			return list, err0
		}
		if len(k) > 0 {
			list = append(list, TextPair{k, v})
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return list, nil
}

// TextToMap - create a map from a []byte
func TextToMap(buf []byte) (map[string]string, error) {
	m := make(map[string]string)
	if len(buf) == 0 {
		return m, nil
	}
	l, err := TextToSlice(buf)
	if err != nil {
		return m, err
	}
	if len(l) > 0 {
		for _, p := range l {
			m[p.Key] = p.Value
		}
	}
	return m, nil
}

func parseLine(line string, isMap bool) (string, string, error) {
	if len(line) == 0 {
		return "", "", nil
	}
	line = strings.TrimLeft(line, " ")
	if isEmpty(line) || isComment(line) {
		return "", "", nil
	}
	var key string
	var val string
	i := strings.Index(line, delimiter)
	if isMap {
		if i == -1 {
			return "", "", fmt.Errorf("invalid argument : line does not contain the ':' delimeter : [%v]", line)
		}
		key = line[:i]
		val = line[i+1:]
		val = removeCrLf(val)
	} else {
		key = removeCrLf(line)
	}
	return strings.TrimSpace(key), strings.TrimLeft(val, " "), nil
}

func isEmpty(line string) bool {
	return len(line) == 0 || line == "" || line == "\r\n" || line == "\n"
}

func isComment(line string) bool {
	return strings.HasPrefix(line, comment)
}

func removeCrLf(s string) string {
	index := strings.Index(s, "\r")
	if index != -1 {
		s = s[:index]
	}
	index = strings.Index(s, "\n")
	if index != -1 {
		s = s[:index]
	}
	return s
}
